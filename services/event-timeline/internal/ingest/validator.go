package ingest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

type EventInput struct {
	StreamID       string            `json:"stream_id"`
	EventType      string            `json:"event_type"`
	Payload        json.RawMessage   `json:"payload"`
	Metadata       map[string]string `json:"metadata"`
	IdempotencyKey string            `json:"idempotency_key"`
	OccurredAt     time.Time         `json:"occurred_at"`
	SchemaVersion  int               `json:"schema_version"`
}

func ValidateEventInput(in EventInput) error {
	if in.StreamID == "" || in.EventType == "" || len(in.Payload) == 0 || in.OccurredAt.IsZero() {
		return fmt.Errorf("missing required fields: %w", domain.ErrValidation)
	}
	return nil
}
