package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
)

// TemporalQuery executes time-range queries
type TemporalQuery struct {
	client *elasticsearch.Client
	logger *slog.Logger
}

// NewTemporalQuery creates a new temporal query
func NewTemporalQuery(client *elasticsearch.Client, logger *slog.Logger) *TemporalQuery {
	return &TemporalQuery{
		client: client,
		logger: logger,
	}
}

// Execute runs the temporal query
func (tq *TemporalQuery) Execute(ctx context.Context, query *domain.TemporalQuery) (*domain.SearchResults, error) {
	filters := []map[string]interface{}{}

	if query.StreamID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"stream_id": query.StreamID},
		})
	}

	q := map[string]interface{}{
		"from": query.Page * query.Size,
		"size": query.Size,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"range": map[string]interface{}{
							"occurred_at": map[string]interface{}{
								"gte": query.From,
								"lte": query.To,
							},
						},
					},
				},
				"filter": filters,
			},
		},
	}

	body, _ := json.Marshal(q)
	res, err := tq.client.Search(tq.client.Search.WithContext(ctx), tq.client.Search.WithIndex("aevum-events"), tq.client.Search.WithBody(bytes.NewBufferString(string(body))))
	if err != nil {
		return nil, fmt.Errorf("temporal query failed: %w", err)
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return parseResults(result, query.Page, query.Size), nil
}

// parseResults parses ES response into SearchResults
func parseResults(esResult map[string]interface{}, page, size int) *domain.SearchResults {
	hits := esResult["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))

	results := &domain.SearchResults{
		Total:    total,
		Hits:     []interface{}{},
		TimeMs:   int64(esResult["took"].(float64)),
		NextPage: page + 1,
		HasMore:  (page+1)*size < int(total),
	}

	if hitList, ok := hits["hits"].([]interface{}); ok {
		for _, h := range hitList {
			if hit, ok := h.(map[string]interface{}); ok {
				results.Hits = append(results.Hits, hit["_source"])
			}
		}
	}

	return results
}

// CorrelationQuery finds correlated events and decisions
type CorrelationQuery struct {
	client *elasticsearch.Client
	logger *slog.Logger
}

// NewCorrelationQuery creates a new correlation query
func NewCorrelationQuery(client *elasticsearch.Client, logger *slog.Logger) *CorrelationQuery {
	return &CorrelationQuery{
		client: client,
		logger: logger,
	}
}

// Execute runs the correlation query
func (cq *CorrelationQuery) Execute(ctx context.Context, query *domain.CorrelationQuery) (*domain.SearchResults, error) {
	filters := []map[string]interface{}{}

	if query.EventID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"event_id": query.EventID},
		})
	}
	if query.DecisionID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"decision_id": query.DecisionID},
		})
	}
	if query.RuleID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"rule_id": query.RuleID},
		})
	}
	if query.StreamID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"stream_id": query.StreamID},
		})
	}

	q := map[string]interface{}{
		"from":  query.Page * query.Size,
		"size":  query.Size,
		"query": map[string]interface{}{"bool": map[string]interface{}{"filter": filters}},
	}

	body, _ := json.Marshal(q)
	res, err := cq.client.Search(cq.client.Search.WithContext(ctx), cq.client.Search.WithIndex("aevum-events", "aevum-decisions"), cq.client.Search.WithBody(bytes.NewBufferString(string(body))))
	if err != nil {
		return nil, fmt.Errorf("correlation query failed: %w", err)
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return parseResults(result, query.Page, query.Size), nil
}
