package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/elastic/go-elasticsearch/v8"
)

// IndexManager manages ES indexes
type IndexManager struct {
	client *elasticsearch.Client
}

// NewIndexManager creates a new index manager
func NewIndexManager(esClient *elasticsearch.Client) *IndexManager {
	return &IndexManager{client: esClient}
}

// CreateIndexes creates all required indexes
func (im *IndexManager) CreateIndexes(ctx context.Context) error {
	if err := im.createIndex(ctx, "aevum-events", EventMapping); err != nil {
		return err
	}
	if err := im.createIndex(ctx, "aevum-decisions", DecisionMapping); err != nil {
		return err
	}
	if err := im.createIndex(ctx, "aevum-sync-state", SyncStateMapping); err != nil {
		return err
	}
	return nil
}

// createIndex creates a single index
func (im *IndexManager) createIndex(ctx context.Context, indexName string, mapping string) error {
	exists, err := im.IndexExists(ctx, indexName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	res, err := im.client.Indices.Create(indexName, im.client.Indices.Create.WithContext(ctx), im.client.Indices.Create.WithBody(bytes.NewBufferString(mapping)))
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to create index %s: %s", indexName, string(body))
	}
	return nil
}

// IndexExists checks if an index exists
func (im *IndexManager) IndexExists(ctx context.Context, indexName string) (bool, error) {
	res, err := im.client.Indices.Exists([]string{indexName}, im.client.Indices.Exists.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	return res.StatusCode == 200, nil
}

// DeleteIndex deletes an index
func (im *IndexManager) DeleteIndex(ctx context.Context, indexName string) error {
	res, err := im.client.Indices.Delete([]string{indexName}, im.client.Indices.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// GetIndexStats returns index statistics
func (im *IndexManager) GetIndexStats(ctx context.Context, indexName string) (map[string]interface{}, error) {
	res, err := im.client.Indices.Stats(im.client.Indices.Stats.WithIndex(indexName), im.client.Indices.Stats.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("failed to get stats: status %d", res.StatusCode)
	}
	return map[string]interface{}{"status": "ok"}, nil
}
