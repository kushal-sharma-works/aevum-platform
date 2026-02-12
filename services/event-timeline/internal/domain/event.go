package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	EventID        string            `json:"event_id" dynamodbav:"PK"`
	SK             string            `json:"-" dynamodbav:"SK"`
	StreamID       string            `json:"stream_id" dynamodbav:"GSI1PK"`
	SequenceNumber int64             `json:"sequence_number" dynamodbav:"GSI1SK"`
	EventType      string            `json:"event_type"`
	Payload        json.RawMessage   `json:"payload"`
	Metadata       map[string]string `json:"metadata"`
	IdempotencyKey string            `json:"idempotency_key" dynamodbav:"GSI2PK"`
	OccurredAt     time.Time         `json:"occurred_at"`
	IngestedAt     time.Time         `json:"ingested_at"`
	SchemaVersion  int               `json:"schema_version"`
}

type NewEventInput struct {
	EventID        string
	StreamID       string
	SequenceNumber int64
	EventType      string
	Payload        json.RawMessage
	Metadata       map[string]string
	IdempotencyKey string
	OccurredAt     time.Time
	IngestedAt     time.Time
	SchemaVersion  int
}

func NewEvent(in NewEventInput) (Event, error) {
	if in.EventID == "" || in.StreamID == "" || in.EventType == "" || len(in.Payload) == 0 || in.OccurredAt.IsZero() {
		return Event{}, fmt.Errorf("invalid event input: %w", ErrValidation)
	}
	if in.SchemaVersion <= 0 {
		in.SchemaVersion = 1
	}
	return Event{
		EventID:        in.EventID,
		SK:             fmt.Sprintf("EVENT#%s#%020d", in.StreamID, in.SequenceNumber),
		StreamID:       in.StreamID,
		SequenceNumber: in.SequenceNumber,
		EventType:      in.EventType,
		Payload:        in.Payload,
		Metadata:       in.Metadata,
		IdempotencyKey: in.IdempotencyKey,
		OccurredAt:     in.OccurredAt.UTC(),
		IngestedAt:     in.IngestedAt.UTC(),
		SchemaVersion:  in.SchemaVersion,
	}, nil
}
