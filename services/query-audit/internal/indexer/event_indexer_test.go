package indexer

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"
)

func TestBulkIndexerBufferingWithoutFlush(t *testing.T) {
	bi := NewBulkIndexer(nil, 10, slog.New(slog.NewTextHandler(io.Discard, nil)))
	err := bi.IndexDocument(context.Background(), "aevum-events", "evt-1", map[string]interface{}{"x": 1})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(bi.buffer) != 2 {
		t.Fatalf("expected 2 buffered items, got %d", len(bi.buffer))
	}

	bi.buffer = nil
	err = bi.Flush(context.Background())
	if err != nil {
		t.Fatalf("expected nil flush on empty buffer, got %v", err)
	}
}

func TestConvertEventToIndexed(t *testing.T) {
	input := map[string]interface{}{
		"event_id":        "evt-1",
		"stream_id":       "s-1",
		"sequence_number": float64(7),
		"event_type":      "created",
		"payload":         map[string]interface{}{"amount": 10},
		"metadata":        map[string]interface{}{"source": "test"},
		"occurred_at":     "2026-01-01T00:00:00Z",
		"ingested_at":     "2026-01-01T00:00:01Z",
		"schema_version":  "1",
	}

	indexed := convertEventToIndexed(input)
	if indexed.EventID != "evt-1" || indexed.StreamID != "s-1" {
		t.Fatal("unexpected event conversion")
	}
	if indexed.SequenceNum != 7 {
		t.Fatalf("expected sequence 7, got %d", indexed.SequenceNum)
	}
	if indexed.OccurredAt.IsZero() || indexed.IngestedAt.IsZero() {
		t.Fatal("expected parsed timestamps")
	}
}

func TestConvertDecisionToIndexedAndHelpers(t *testing.T) {
	input := map[string]interface{}{
		"decision_id":        "dec-1",
		"event_id":           "evt-1",
		"stream_id":          "s-1",
		"rule_id":            "rule-1",
		"rule_version":       "2",
		"status":             "approved",
		"deterministic_hash": "hash-1",
		"input":              map[string]interface{}{"x": 1},
		"output":             map[string]interface{}{"approved": true},
		"trace": []interface{}{
			map[string]interface{}{
				"step":      float64(1),
				"condition": "amount > 10",
				"result":    true,
				"message":   "ok",
				"timestamp": "2026-01-01T00:00:01Z",
			},
		},
	}

	indexed := convertDecisionToIndexed(input)
	if indexed.DecisionID != "dec-1" {
		t.Fatal("unexpected decision id")
	}
	if len(indexed.Trace) != 1 || !indexed.Trace[0].Result {
		t.Fatal("unexpected trace conversion")
	}

	if toInt64(int64(4)) != 4 || toInt64(float64(5)) != 5 || toInt64("x") != 0 {
		t.Fatal("unexpected toInt64 behavior")
	}
	if !toBool(true) || toBool("true") {
		t.Fatal("unexpected toBool behavior")
	}
	if toString("abc") != "abc" || toString(1) != "" {
		t.Fatal("unexpected toString behavior")
	}
	if !parseTime("2026-01-01T00:00:00Z").Equal(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("unexpected parseTime behavior")
	}
	if parseTime(10).IsZero() {
		t.Fatal("expected non-zero parseTime fallback")
	}
}
