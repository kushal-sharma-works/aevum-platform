package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

type DynamoDBEventStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBEventStore(client *dynamodb.Client, tableName string) *DynamoDBEventStore {
	return &DynamoDBEventStore{client: client, tableName: tableName}
}

const (
	batchWriteMaxItems = 25
	batchWriteMaxRetry = 5
)

func idempotencyLookupKey(streamID, key string) string {
	if key == "" {
		return ""
	}
	return streamID + "#" + key
}

func sequenceGuardPK(streamID string) string {
	return "SEQ#" + streamID
}

func sequenceGuardSK(sequence int64) string {
	return strconv.FormatInt(sequence, 10)
}

func (s *DynamoDBEventStore) PutEvent(ctx context.Context, event domain.Event) error {
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	transactItems := []types.TransactWriteItem{
		{
			Put: &types.Put{
				TableName:           aws.String(s.tableName),
				Item:                item,
				ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
			},
		},
		{
			Put: &types.Put{
				TableName: aws.String(s.tableName),
				Item: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: sequenceGuardPK(event.StreamID)},
					"SK": &types.AttributeValueMemberS{Value: sequenceGuardSK(event.SequenceNumber)},
				},
				ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
			},
		},
	}

	idempotencyIndex := -1
	if event.IdempotencyKey != "" {
		idempotencyIndex = len(transactItems)
		transactItems = append(transactItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName: aws.String(s.tableName),
				Item: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: "IDEMP#" + idempotencyLookupKey(event.StreamID, event.IdempotencyKey)},
					"SK": &types.AttributeValueMemberS{Value: "LOCK"},
				},
				ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
			},
		})
	}

	_, err = s.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{TransactItems: transactItems})
	if err != nil {
		var cancelled *types.TransactionCanceledException
		if errors.As(err, &cancelled) {
			if idempotencyIndex >= 0 && len(cancelled.CancellationReasons) > idempotencyIndex {
				if aws.ToString(cancelled.CancellationReasons[idempotencyIndex].Code) == "ConditionalCheckFailed" {
					return fmt.Errorf("idempotency conflict: %w", domain.ErrIdempotencyConflict)
				}
			}
			for _, reason := range cancelled.CancellationReasons {
				if aws.ToString(reason.Code) == "ConditionalCheckFailed" {
					return fmt.Errorf("sequence conflict: %w", domain.ErrSequenceConflict)
				}
			}
		}
		var ccf *types.ConditionalCheckFailedException
		if errors.As(err, &ccf) {
			return fmt.Errorf("sequence conflict: %w", domain.ErrSequenceConflict)
		}
		return fmt.Errorf("put event: %w", err)
	}
	return nil
}

func (s *DynamoDBEventStore) PutEventsBatch(ctx context.Context, events []domain.Event) error {
	if len(events) == 0 {
		return nil
	}

	for start := 0; start < len(events); start += batchWriteMaxItems {
		end := start + batchWriteMaxItems
		if end > len(events) {
			end = len(events)
		}

		chunk := events[start:end]
		items := make([]types.WriteRequest, 0, len(chunk))
		for _, event := range chunk {
			item, err := attributevalue.MarshalMap(event)
			if err != nil {
				return fmt.Errorf("marshal event: %w", err)
			}
			items = append(items, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
		}

		pending := map[string][]types.WriteRequest{s.tableName: items}
		for attempt := 0; attempt < batchWriteMaxRetry && len(pending[s.tableName]) > 0; attempt++ {
			resp, err := s.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{RequestItems: pending})
			if err != nil {
				return fmt.Errorf("batch write item: %w", err)
			}
			pending = resp.UnprocessedItems
			if len(pending[s.tableName]) == 0 {
				break
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("batch write cancelled: %w", ctx.Err())
			case <-time.After(time.Duration(attempt+1) * 25 * time.Millisecond):
			}
		}
		if len(pending[s.tableName]) > 0 {
			return fmt.Errorf("batch write item: unprocessed items remain")
		}
	}
	return nil
}

func (s *DynamoDBEventStore) GetByEventID(ctx context.Context, eventID string) (domain.Event, error) {
	resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		KeyConditionExpression: aws.String("PK = :event_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":event_id": &types.AttributeValueMemberS{Value: eventID},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return domain.Event{}, fmt.Errorf("get event: %w", err)
	}
	if len(resp.Items) == 0 {
		return domain.Event{}, fmt.Errorf("event not found: %w", domain.ErrNotFound)
	}
	var event domain.Event
	if err := attributevalue.UnmarshalMap(resp.Items[0], &event); err != nil {
		return domain.Event{}, fmt.Errorf("unmarshal event: %w", err)
	}
	return event, nil
}

func (s *DynamoDBEventStore) FindByIdempotencyKey(ctx context.Context, streamID, key string) (domain.Event, error) {
	resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String(GSI2Name),
		KeyConditionExpression: aws.String("GSI2PK = :idempotency_key"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":idempotency_key": &types.AttributeValueMemberS{Value: idempotencyLookupKey(streamID, key)},
		},
		Limit: aws.Int32(1),
	})
	if err != nil {
		return domain.Event{}, fmt.Errorf("query by idempotency key: %w", err)
	}
	if len(resp.Items) == 0 {
		return domain.Event{}, fmt.Errorf("idempotency key not found: %w", domain.ErrNotFound)
	}
	var event domain.Event
	if err := attributevalue.UnmarshalMap(resp.Items[0], &event); err != nil {
		return domain.Event{}, fmt.Errorf("unmarshal event: %w", err)
	}
	return event, nil
}

func (s *DynamoDBEventStore) GetLatestSequence(ctx context.Context, streamID string) (int64, error) {
	resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String(GSI1Name),
		KeyConditionExpression: aws.String("GSI1PK = :stream_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":stream_id": &types.AttributeValueMemberS{Value: streamID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	})
	if err != nil {
		return 0, fmt.Errorf("query latest sequence: %w", err)
	}
	if len(resp.Items) == 0 {
		return 0, nil
	}
	var event domain.Event
	if err := attributevalue.UnmarshalMap(resp.Items[0], &event); err != nil {
		return 0, fmt.Errorf("unmarshal event: %w", err)
	}
	return event.SequenceNumber, nil
}

func (s *DynamoDBEventStore) QueryByStream(ctx context.Context, streamID string, fromSequence int64, direction string, limit int32) ([]domain.Event, int64, bool, error) {
	if limit <= 0 {
		limit = 50
	}
	scanForward := direction != domain.DirectionBackward
	operator := ">="
	if !scanForward {
		operator = "<="
	}
	resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String(GSI1Name),
		KeyConditionExpression: aws.String("GSI1PK = :stream_id AND GSI1SK " + operator + " :seq"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":stream_id": &types.AttributeValueMemberS{Value: streamID},
			":seq":       &types.AttributeValueMemberN{Value: strconv.FormatInt(fromSequence, 10)},
		},
		ScanIndexForward: aws.Bool(scanForward),
		Limit:            aws.Int32(limit),
	})
	if err != nil {
		return nil, 0, false, fmt.Errorf("query stream: %w", err)
	}
	events := make([]domain.Event, 0, len(resp.Items))
	for _, item := range resp.Items {
		var event domain.Event
		if err := attributevalue.UnmarshalMap(item, &event); err != nil {
			return nil, 0, false, fmt.Errorf("unmarshal stream event: %w", err)
		}
		events = append(events, event)
	}
	nextSeq := fromSequence
	if len(events) > 0 {
		nextSeq = events[len(events)-1].SequenceNumber
		if scanForward {
			nextSeq++
		} else {
			nextSeq--
		}
	}
	return events, nextSeq, len(resp.LastEvaluatedKey) > 0, nil
}
