package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// BulkIndexer manages bulk document indexing
type BulkIndexer struct {
	client    *elasticsearch.Client
	batchSize int
	buffer    []interface{}
	logger    *slog.Logger
}

// NewBulkIndexer creates a new bulk indexer
func NewBulkIndexer(client *elasticsearch.Client, batchSize int, logger *slog.Logger) *BulkIndexer {
	return &BulkIndexer{
		client:    client,
		batchSize: batchSize,
		buffer:    []interface{}{},
		logger:    logger,
	}
}

// IndexDocument adds a document to the bulk buffer
func (bi *BulkIndexer) IndexDocument(ctx context.Context, indexName, docID string, doc interface{}) error {
	if err := bi.addToBatch(indexName, docID, doc); err != nil {
		return err
	}

	if len(bi.buffer) >= bi.batchSize*2 {
		return bi.Flush(ctx)
	}
	return nil
}

// addToBatch adds a document to the batch
func (bi *BulkIndexer) addToBatch(indexName, docID string, doc interface{}) error {
	metadata := map[string]interface{}{
		"index": map[string]interface{}{
			"_index": indexName,
			"_id":    docID,
		},
	}
	bi.buffer = append(bi.buffer, metadata, doc)
	return nil
}

// Flush sends all buffered documents to ES
func (bi *BulkIndexer) Flush(ctx context.Context) error {
	if len(bi.buffer) == 0 {
		return nil
	}

	var body bytes.Buffer
	for _, item := range bi.buffer {
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal item: %w", err)
		}
		body.Write(data)
		body.WriteString("\n")
	}

	req := esapi.BulkRequest{
		Body: &body,
	}
	res, err := req.Do(ctx, bi.client)
	if err != nil {
		return fmt.Errorf("bulk request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk indexing failed: status %d", res.StatusCode)
	}

	bi.buffer = []interface{}{}
	bi.logger.Info("bulk indexing completed", slog.Int("documents", len(bi.buffer)/2))
	return nil
}
