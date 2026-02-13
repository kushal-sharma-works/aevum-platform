package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
)

// DiffEngine compares decision sets
type DiffEngine struct {
	client *elasticsearch.Client
	logger *slog.Logger
}

// NewDiffEngine creates a new diff engine
func NewDiffEngine(client *elasticsearch.Client, logger *slog.Logger) *DiffEngine {
	return &DiffEngine{
		client: client,
		logger: logger,
	}
}

// Compare compares decisions at two points in time
func (de *DiffEngine) Compare(ctx context.Context, query *domain.DiffQuery) (*domain.DiffResult, error) {
	t1Decisions, err := de.queryDecisions(ctx, query.T1, query.RuleID, query.StreamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query T1 decisions: %w", err)
	}

	t2Decisions, err := de.queryDecisions(ctx, query.T2, query.RuleID, query.StreamID)
	if err != nil {
		return nil, fmt.Errorf("failed to query T2 decisions: %w", err)
	}

	return de.diff(t1Decisions, t2Decisions), nil
}

// queryDecisions fetches decisions at a specific time
func (de *DiffEngine) queryDecisions(ctx context.Context, before time.Time, ruleID, streamID string) (map[string]*domain.IndexedDecision, error) {
	filters := []map[string]interface{}{
		{
			"range": map[string]interface{}{
				"evaluated_at": map[string]interface{}{"lte": before},
			},
		},
	}

	if ruleID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"rule_id": ruleID},
		})
	}

	if streamID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"stream_id": streamID},
		})
	}

	q := map[string]interface{}{
		"size": 10000,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{"filter": filters},
		},
	}

	body, _ := json.Marshal(q)
	res, err := de.client.Search(de.client.Search.WithContext(ctx), de.client.Search.WithIndex("aevum-decisions"), de.client.Search.WithBody(bytes.NewBufferString(string(body))))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	decisions := make(map[string]*domain.IndexedDecision)
	if hits, ok := result["hits"].(map[string]interface{})["hits"].([]interface{}); ok {
		for _, h := range hits {
			if hit, ok := h.(map[string]interface{}); ok {
				if source, ok := hit["_source"].(map[string]interface{}); ok {
					decisionID := fmt.Sprint(source["decision_id"])
					decisions[decisionID] = &domain.IndexedDecision{
						DecisionID: decisionID,
						RuleID:     fmt.Sprint(source["rule_id"]),
					}
				}
			}
		}
	}

	return decisions, nil
}

// diff computes the difference between two decision sets
func (de *DiffEngine) diff(t1, t2 map[string]*domain.IndexedDecision) *domain.DiffResult {
	result := &domain.DiffResult{
		Added:   []string{},
		Removed: []string{},
		Changed: []domain.FieldDiff{},
	}

	// Find removed decisions
	for id := range t1 {
		if _, exists := t2[id]; !exists {
			result.Removed = append(result.Removed, id)
		}
	}

	// Find added decisions
	for id := range t2 {
		if _, exists := t1[id]; !exists {
			result.Added = append(result.Added, id)
		}
	}

	result.Summary = fmt.Sprintf("Added: %d, Removed: %d", len(result.Added), len(result.Removed))
	return result
}
