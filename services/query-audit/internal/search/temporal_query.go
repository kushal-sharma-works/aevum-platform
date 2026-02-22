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
		"from": (query.Page - 1) * query.Size,
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
	indexes := temporalIndexes(query.Type)
	res, err := tq.client.Search(tq.client.Search.WithContext(ctx), tq.client.Search.WithIndex(indexes...), tq.client.Search.WithBody(bytes.NewBufferString(string(body))))
	if err != nil {
		return nil, fmt.Errorf("temporal query failed: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("temporal query failed: status %d", res.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return parseResults(result, query.Page, query.Size), nil
}

func temporalIndexes(queryType string) []string {
	switch queryType {
	case "events":
		return []string{"aevum-events"}
	case "decisions":
		return []string{"aevum-decisions"}
	default:
		return []string{"aevum-events", "aevum-decisions"}
	}
}

// parseResults parses ES response into SearchResults
func parseResults(esResult map[string]interface{}, page, size int) *domain.SearchResults {
	hitsContainer, ok := esResult["hits"].(map[string]interface{})
	if !ok {
		return &domain.SearchResults{Hits: []interface{}{}, NextPage: page + 1}
	}

	total := int64(0)
	if totalContainer, ok := hitsContainer["total"].(map[string]interface{}); ok {
		if value, ok := totalContainer["value"].(float64); ok {
			total = int64(value)
		}
	}

	took := int64(0)
	if tookValue, ok := esResult["took"].(float64); ok {
		took = int64(tookValue)
	}

	results := &domain.SearchResults{
		Total:    total,
		Hits:     []interface{}{},
		TimeMs:   took,
		NextPage: page + 1,
		HasMore:  (page+1)*size < int(total),
	}

	if hitList, ok := hitsContainer["hits"].([]interface{}); ok {
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
	if query.RuleVersion != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"rule_version": query.RuleVersion},
		})
	}
	if query.StreamID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"stream_id": query.StreamID},
		})
	}
	if query.EventType != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"event_type": query.EventType},
		})
	}

	q := map[string]interface{}{
		"from":  (query.Page - 1) * query.Size,
		"size":  query.Size,
		"query": map[string]interface{}{"bool": map[string]interface{}{"filter": filters}},
	}

	body, _ := json.Marshal(q)
	res, err := cq.client.Search(cq.client.Search.WithContext(ctx), cq.client.Search.WithIndex("aevum-events", "aevum-decisions"), cq.client.Search.WithBody(bytes.NewBufferString(string(body))))
	if err != nil {
		return nil, fmt.Errorf("correlation query failed: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("correlation query failed: status %d", res.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return parseResults(result, query.Page, query.Size), nil
}
