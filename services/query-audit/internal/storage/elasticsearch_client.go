package storage

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
)

// ElasticsearchClient wraps the ES client
type ElasticsearchClient struct {
	client *elasticsearch.Client
}

// NewElasticsearchClient creates a new ES client
func NewElasticsearchClient(urls []string) (*ElasticsearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: urls,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}
	return &ElasticsearchClient{client: client}, nil
}

// Health checks the ES cluster health
func (ec *ElasticsearchClient) Health(ctx context.Context) error {
	res, err := ec.client.Cluster.Health(ec.client.Cluster.Health.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to check cluster health: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("cluster unhealthy: status %d", res.StatusCode)
	}
	return nil
}

// GetClient returns the underlying ES client
func (ec *ElasticsearchClient) GetClient() *elasticsearch.Client {
	return ec.client
}
