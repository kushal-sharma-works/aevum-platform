package indexer

import (
	"context"
	"log/slog"
	"time"

	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/clients"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/domain"
)

// EventIndexer synchronizes and indexes events
type EventIndexer struct {
	client      *clients.EventTimelineClient
	bulkIndexer *BulkIndexer
	logger      *slog.Logger
}

// NewEventIndexer creates a new event indexer
func NewEventIndexer(client *clients.EventTimelineClient, bulkIndexer *BulkIndexer, logger *slog.Logger) *EventIndexer {
	return &EventIndexer{
		client:      client,
		bulkIndexer: bulkIndexer,
		logger:      logger,
	}
}

// Sync fetches events and indexes them
func (ei *EventIndexer) Sync(ctx context.Context, cursor string) (string, error) {
	events, newCursor, err := ei.client.GetEvents(ctx, cursor, 100)
	if err != nil {
		ei.logger.Error("failed to fetch events", slog.Any("error", err))
		return cursor, err
	}

	for _, event := range events {
		indexed := convertEventToIndexed(event)
		if err := ei.bulkIndexer.IndexDocument(ctx, "aevum-events", indexed.EventID, indexed); err != nil {
			ei.logger.Error("failed to index event", slog.String("event_id", indexed.EventID), slog.Any("error", err))
		}
	}

	if err := ei.bulkIndexer.Flush(ctx); err != nil {
		ei.logger.Error("failed to flush bulk indexer", slog.Any("error", err))
		return cursor, err
	}

	return newCursor, nil
}

// convertEventToIndexed transforms a raw event to IndexedEvent
func convertEventToIndexed(event map[string]interface{}) *domain.IndexedEvent {
	payload := convertToStringMap(event["payload"])
	metadata := convertToStringMap(event["metadata"])

	occurredAt := parseTime(event["occurred_at"])
	ingestedAt := parseTime(event["ingested_at"])

	return &domain.IndexedEvent{
		EventID:       toString(event["event_id"]),
		StreamID:      toString(event["stream_id"]),
		SequenceNum:   toInt64(event["sequence_number"]),
		EventType:     toString(event["event_type"]),
		Payload:       payload,
		Metadata:      metadata,
		OccurredAt:    occurredAt,
		IngestedAt:    ingestedAt,
		SchemaVersion: toString(event["schema_version"]),
	}
}

// convertToStringMap safely converts to map[string]interface{}
func convertToStringMap(v interface{}) map[string]interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		return m
	}
	return make(map[string]interface{})
}

// parseTime parses ISO8601 timestamp
func parseTime(v interface{}) time.Time {
	if s, ok := v.(string); ok {
		t, _ := time.Parse(time.RFC3339, s)
		return t
	}
	return time.Now()
}

// toString converts to string
func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

// toInt64 converts to int64
func toInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case float64:
		return int64(val)
	default:
		return 0
	}
}

// DecisionIndexer synchronizes and indexes decisions
type DecisionIndexer struct {
	client      *clients.DecisionEngineClient
	bulkIndexer *BulkIndexer
	logger      *slog.Logger
}

// NewDecisionIndexer creates a new decision indexer
func NewDecisionIndexer(client *clients.DecisionEngineClient, bulkIndexer *BulkIndexer, logger *slog.Logger) *DecisionIndexer {
	return &DecisionIndexer{
		client:      client,
		bulkIndexer: bulkIndexer,
		logger:      logger,
	}
}

// Sync fetches decisions and indexes them
func (di *DecisionIndexer) Sync(ctx context.Context, from, to time.Time) error {
	decisions, err := di.client.GetDecisions(ctx, from, to)
	if err != nil {
		di.logger.Error("failed to fetch decisions", slog.Any("error", err))
		return err
	}

	for _, decision := range decisions {
		indexed := convertDecisionToIndexed(decision)
		if err := di.bulkIndexer.IndexDocument(ctx, "aevum-decisions", indexed.DecisionID, indexed); err != nil {
			di.logger.Error("failed to index decision", slog.String("decision_id", indexed.DecisionID), slog.Any("error", err))
		}
	}

	if err := di.bulkIndexer.Flush(ctx); err != nil {
		di.logger.Error("failed to flush bulk indexer", slog.Any("error", err))
		return err
	}

	return nil
}

// convertDecisionToIndexed transforms a raw decision to IndexedDecision
func convertDecisionToIndexed(decision map[string]interface{}) *domain.IndexedDecision {
	input := convertToStringMap(decision["input"])
	output := convertToStringMap(decision["output"])

	trace := []domain.TraceEntry{}
	if traceList, ok := decision["trace"].([]interface{}); ok {
		for _, t := range traceList {
			if traceMap, ok := t.(map[string]interface{}); ok {
				trace = append(trace, domain.TraceEntry{
					Step:      int(toInt64(traceMap["step"])),
					Condition: toString(traceMap["condition"]),
					Result:    toBool(traceMap["result"]),
					Message:   toString(traceMap["message"]),
					Timestamp: parseTime(traceMap["timestamp"]),
				})
			}
		}
	}

	return &domain.IndexedDecision{
		DecisionID:        toString(decision["decision_id"]),
		EventID:           toString(decision["event_id"]),
		StreamID:          toString(decision["stream_id"]),
		RuleID:            toString(decision["rule_id"]),
		RuleVersion:       toString(decision["rule_version"]),
		Status:            toString(decision["status"]),
		DeterministicHash: toString(decision["deterministic_hash"]),
		Input:             input,
		Output:            output,
		Trace:             trace,
		EvaluatedAt:       parseTime(decision["evaluated_at"]),
		EventOccurredAt:   parseTime(decision["event_occurred_at"]),
	}
}

// toBool converts to bool
func toBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}
