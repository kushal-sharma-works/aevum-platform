package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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

func (s *DynamoDBEventStore) PutEvent(ctx context.Context, event domain.Event) error {
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(s.tableName),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(PK) AND attribute_not_exists(SK)"),
	})
	if err != nil {
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
	items := make([]types.WriteRequest, 0, len(events))
	for _, event := range events {
		item, err := attributevalue.MarshalMap(event)
		if err != nil {
			return fmt.Errorf("marshal event: %w", err)
		}
		items = append(items, types.WriteRequest{PutRequest: &types.PutRequest{Item: item}})
	}
	_, err := s.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{s.tableName: items},
	})
	if err != nil {
		return fmt.Errorf("batch write item: %w", err)
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

func (s *DynamoDBEventStore) FindByIdempotencyKey(ctx context.Context, key string) (domain.Event, error) {
	resp, err := s.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		IndexName:              aws.String(GSI2Name),
		KeyConditionExpression: aws.String("GSI2PK = :idempotency_key"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":idempotency_key": &types.AttributeValueMemberS{Value: key},
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
