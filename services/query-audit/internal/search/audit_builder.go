package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/clients"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
)

// AuditBuilder builds complete causal chains
type AuditBuilder struct {
	esClient       *elasticsearch.Client
	eventClient    *clients.EventTimelineClient
	decisionClient *clients.DecisionEngineClient
	logger         *slog.Logger
}

// NewAuditBuilder creates a new audit builder
func NewAuditBuilder(esClient *elasticsearch.Client, eventClient *clients.EventTimelineClient, decisionClient *clients.DecisionEngineClient, logger *slog.Logger) *AuditBuilder {
	return &AuditBuilder{
		esClient:       esClient,
		eventClient:    eventClient,
		decisionClient: decisionClient,
		logger:         logger,
	}
}

// Build builds an audit trail for a decision
func (ab *AuditBuilder) Build(ctx context.Context, decisionID string) (*domain.AuditTrail, error) {
	// Fetch decision from ES
	decision, err := ab.fetchDecision(ctx, decisionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch decision: %w", err)
	}

	// Fetch event from Event Timeline Service
	eventRaw, err := ab.eventClient.GetEvent(ctx, decision.EventID)
	if err != nil {
		ab.logger.Warn("failed to fetch event", slog.String("event_id", decision.EventID), slog.Any("error", err))
	}

	// Fetch rule definition from Decision Engine Service
	rule, err := ab.decisionClient.GetRule(ctx, decision.RuleID)
	if err != nil {
		ab.logger.Warn("failed to fetch rule", slog.String("rule_id", decision.RuleID), slog.Any("error", err))
	}

	// Build chain - pass raw event and convert in function
	chain := ab.buildChain(decision, eventRaw, rule)

	// Convert event map to IndexedEvent if present
	var event *domain.IndexedEvent
	if eventRaw != nil {
		event = &domain.IndexedEvent{
			EventID:  fmt.Sprint((*eventRaw)["event_id"]),
			StreamID: fmt.Sprint((*eventRaw)["stream_id"]),
		}
	}

	return &domain.AuditTrail{
		Decision:       decision,
		Event:          event,
		RuleDefinition: rule,
		Chain:          chain,
	}, nil
}

// fetchDecision fetches a decision from ES
func (ab *AuditBuilder) fetchDecision(ctx context.Context, decisionID string) (*domain.IndexedDecision, error) {
	res, err := ab.esClient.Get("aevum-decisions", decisionID, ab.esClient.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("fetch decision failed: status %d", res.StatusCode)
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	source, ok := payload["_source"].(map[string]interface{})
	if !ok {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("decision source missing: %s", string(body))
	}

	return &domain.IndexedDecision{
		DecisionID:  asString(source["decision_id"]),
		EventID:     asString(source["event_id"]),
		StreamID:    asString(source["stream_id"]),
		RuleID:      asString(source["rule_id"]),
		RuleVersion: asString(source["rule_version"]),
		Status:      asString(source["status"]),
		EvaluatedAt: asTime(source["evaluated_at"]),
	}, nil
}

// buildChain builds the causal chain
func (ab *AuditBuilder) buildChain(decision *domain.IndexedDecision, eventRaw *map[string]interface{}, rule interface{}) []domain.AuditStep {
	chain := []domain.AuditStep{}

	// Event step
	if eventRaw != nil {
		chain = append(chain, domain.AuditStep{
			Type:        "event_occurred",
			Description: "Event occurred in stream",
			Data: map[string]interface{}{
				"event_id":   (*eventRaw)["event_id"],
				"stream_id":  (*eventRaw)["stream_id"],
				"event_type": (*eventRaw)["event_type"],
			},
			Timestamp: time.Now(),
		})
	}

	// Decision evaluation step
	chain = append(chain, domain.AuditStep{
		Type:        "decision_evaluated",
		Description: "Decision rule evaluated",
		Data: map[string]interface{}{
			"rule_id":      decision.RuleID,
			"rule_version": decision.RuleVersion,
			"status":       decision.Status,
		},
		Timestamp: decision.EvaluatedAt,
	})

	// Trace steps
	for _, trace := range decision.Trace {
		chain = append(chain, domain.AuditStep{
			Type:        "condition_evaluated",
			Description: trace.Condition,
			Data: map[string]interface{}{
				"condition": trace.Condition,
				"result":    trace.Result,
				"message":   trace.Message,
			},
			Timestamp: trace.Timestamp,
		})
	}

	return chain
}

func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

func asTime(v interface{}) time.Time {
	if s, ok := v.(string); ok {
		if parsed, err := time.Parse(time.RFC3339, s); err == nil {
			return parsed
		}
	}
	return time.Time{}
}
