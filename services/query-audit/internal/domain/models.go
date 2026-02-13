package domain

import "time"

// IndexedEvent represents an event in Elasticsearch
type IndexedEvent struct {
	EventID       string                 `json:"event_id"`
	StreamID      string                 `json:"stream_id"`
	SequenceNum   int64                  `json:"sequence_number"`
	EventType     string                 `json:"event_type"`
	Payload       map[string]interface{} `json:"payload"`
	Metadata      map[string]interface{} `json:"metadata"`
	OccurredAt    time.Time              `json:"occurred_at"`
	IngestedAt    time.Time              `json:"ingested_at"`
	SchemaVersion string                 `json:"schema_version"`
}

// IndexedDecision represents a decision in Elasticsearch
type IndexedDecision struct {
	DecisionID        string                 `json:"decision_id"`
	EventID           string                 `json:"event_id"`
	StreamID          string                 `json:"stream_id"`
	RuleID            string                 `json:"rule_id"`
	RuleVersion       string                 `json:"rule_version"`
	Status            string                 `json:"status"`
	DeterministicHash string                 `json:"deterministic_hash"`
	Input             map[string]interface{} `json:"input"`
	Output            map[string]interface{} `json:"output"`
	Trace             []TraceEntry           `json:"trace"`
	EvaluatedAt       time.Time              `json:"evaluated_at"`
	EventOccurredAt   time.Time              `json:"event_occurred_at"`
}

// TraceEntry represents a single step in the decision trace
type TraceEntry struct {
	Step      int       `json:"step"`
	Condition string    `json:"condition"`
	Result    bool      `json:"result"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// AuditTrail represents the complete causal chain
type AuditTrail struct {
	Decision       *IndexedDecision `json:"decision"`
	Event          *IndexedEvent    `json:"event"`
	RuleDefinition interface{}      `json:"rule_definition"`
	Chain          []AuditStep      `json:"chain"`
}

// AuditStep is a single step in the audit chain
type AuditStep struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   time.Time              `json:"timestamp"`
}

// DiffResult contains differences between two decision sets
type DiffResult struct {
	Added   []string    `json:"added"`
	Removed []string    `json:"removed"`
	Changed []FieldDiff `json:"changed"`
	Summary string      `json:"summary"`
}

// FieldDiff represents a change in a specific field
type FieldDiff struct {
	DecisionID string      `json:"decision_id"`
	Field      string      `json:"field"`
	OldValue   interface{} `json:"old_value"`
	NewValue   interface{} `json:"new_value"`
}

// SearchQuery represents a full-text search request
type SearchQuery struct {
	Query    string
	Type     string
	StreamID string
	Page     int
	Size     int
}

// TemporalQuery represents a time-range query
type TemporalQuery struct {
	From     time.Time
	To       time.Time
	StreamID string
	Type     string
	Page     int
	Size     int
}

// CorrelationQuery finds related events and decisions
type CorrelationQuery struct {
	EventID      string
	DecisionID   string
	RuleID       string
	RuleVersion  string
	StreamID     string
	EventType    string
	IncludeAudit bool
	Page         int
	Size         int
}

// DiffQuery compares decisions at two points in time
type DiffQuery struct {
	T1          time.Time
	T2          time.Time
	RuleID      string
	RuleVersion string
	StreamID    string
	Page        int
	Size        int
}

// SearchResults wraps search results with metadata
type SearchResults struct {
	Total    int64         `json:"total"`
	Hits     []interface{} `json:"hits"`
	TimeMs   int64         `json:"time_ms"`
	HasMore  bool          `json:"has_more"`
	NextPage int           `json:"next_page"`
}
