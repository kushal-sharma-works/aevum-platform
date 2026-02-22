package unit

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

func TestNewEventBuildsExpectedRecord(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	ingested := time.Date(2026, 1, 1, 12, 0, 1, 0, time.UTC)

	event, err := domain.NewEvent(domain.NewEventInput{
		EventID:        "evt-1",
		StreamID:       "stream-1",
		SequenceNumber: 7,
		EventType:      "RULE_EVALUATED",
		Payload:        json.RawMessage(`{"k":"v"}`),
		Metadata:       map[string]string{"source": "test"},
		IdempotencyKey: "idem-1",
		OccurredAt:     now,
		IngestedAt:     ingested,
	})

	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}

	if event.SK != "EVENT#stream-1#00000000000000000007" {
		t.Fatalf("unexpected SK: %s", event.SK)
	}

	if event.SchemaVersion != 1 {
		t.Fatalf("expected default schema version 1, got %d", event.SchemaVersion)
	}

	if event.OccurredAt.Location() != time.UTC || event.IngestedAt.Location() != time.UTC {
		t.Fatal("expected UTC timestamps")
	}
}

func TestNewEventRejectsInvalidInput(t *testing.T) {
	_, err := domain.NewEvent(domain.NewEventInput{EventID: "", StreamID: "stream-1"})

	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEventPreservesMetadata(t *testing.T) {
	metadata := map[string]string{"source": "api", "version": "v1"}
	payload := json.RawMessage(`{"action":"test"}`)

	event, _ := domain.NewEvent(domain.NewEventInput{
		EventID:    "evt-2",
		StreamID:   "stream-2",
		EventType:  "TEST",
		Metadata:   metadata,
		Payload:    payload,
		OccurredAt: time.Now().UTC(),
		IngestedAt: time.Now().UTC(),
	})

	if event.Metadata["source"] != "api" || event.Metadata["version"] != "v1" {
		t.Fatal("metadata not preserved")
	}

	if string(event.Payload) != `{"action":"test"}` {
		t.Fatalf("payload corrupted: %s", string(event.Payload))
	}
}
