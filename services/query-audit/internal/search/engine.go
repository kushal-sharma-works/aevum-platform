package search

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8"
)

// Engine performs full-text search
type Engine struct {
	client *elasticsearch.Client
	logger *slog.Logger
}

// NewEngine creates a new search engine
func NewEngine(client *elasticsearch.Client, logger *slog.Logger) *Engine {
	return &Engine{
		client: client,
		logger: logger,
	}
}

// Search performs full-text search
func (e *Engine) Search(ctx context.Context, query string, searchType string, streamID string, from, size int) (map[string]interface{}, error) {
	body := buildSearchQuery(query, searchType, streamID, from, size)

	res, err := e.client.Search(e.client.Search.WithContext(ctx), e.client.Search.WithIndex("aevum-events", "aevum-decisions"), e.client.Search.WithBody(bytes.NewBufferString(body)))
	if err != nil {
		e.logger.Error("search failed", slog.Any("error", err))
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		e.logger.Error("search error", slog.Int("status", res.StatusCode))
		return nil, err
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		e.logger.Error("failed to decode response", slog.Any("error", err))
		return nil, err
	}

	return result, nil
}

// buildSearchQuery builds the ES query
func buildSearchQuery(query string, searchType string, streamID string, from, size int) string {
	filters := []map[string]interface{}{}

	if streamID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"stream_id": streamID,
			},
		})
	}

	if searchType != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"event_type": searchType,
			},
		})
	}

	q := map[string]interface{}{
		"from": from,
		"size": size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":  query,
							"fields": []string{"payload", "metadata", "input", "output"},
						},
					},
				},
				"filter": filters,
			},
		},
	}

	body, _ := json.Marshal(q)
	return string(body)
}
