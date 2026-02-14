package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

type DynamoDBStreamStore struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBStreamStore(client *dynamodb.Client, tableName string) *DynamoDBStreamStore {
	return &DynamoDBStreamStore{client: client, tableName: tableName}
}

func (s *DynamoDBStreamStore) ListStreams(ctx context.Context, limit int32) ([]domain.Stream, error) {
	if limit <= 0 {
		limit = 200
	}
	resp, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:            aws.String(s.tableName),
		Limit:                aws.Int32(limit),
		ProjectionExpression: aws.String("GSI1PK, GSI1SK"),
	})
	if err != nil {
		return nil, fmt.Errorf("scan streams: %w", err)
	}
	seen := map[string]int64{}
	for _, item := range resp.Items {
		streamAttr, ok := item["GSI1PK"].(*types.AttributeValueMemberS)
		if !ok {
			continue
		}
		seqAttr, ok := item["GSI1SK"].(*types.AttributeValueMemberN)
		if !ok {
			continue
		}
		var seq int64
		_, _ = fmt.Sscan(seqAttr.Value, &seq)
		if seq > seen[streamAttr.Value] {
			seen[streamAttr.Value] = seq
		}
	}
	streams := make([]domain.Stream, 0, len(seen))
	for id, latest := range seen {
		streams = append(streams, domain.Stream{StreamID: id, LatestSequence: latest})
	}
	return streams, nil
}
