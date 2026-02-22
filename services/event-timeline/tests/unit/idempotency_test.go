package unit

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
)

func TestIdempotencyFindExisting(t *testing.T) {
	store := &mockEventStore{}
	event, err := domain.NewEvent(domain.NewEventInput{
		EventID:       "id-1",
		StreamID:      "stream-1",
		EventType:     "created",
		Payload:       json.RawMessage(`{"ok":true}`),
		OccurredAt:    time.Now().UTC(),
		IngestedAt:    time.Now().UTC(),
		SchemaVersion: 1,
	})
	require.NoError(t, err)
	store.byIdem = map[string]domain.Event{"stream-1#idem-1": event}

	checker := ingest.NewIdempotencyChecker(store)
	got, ok, err := checker.FindExisting(context.Background(), "stream-1", "idem-1")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, event.EventID, got.EventID)

	_, ok, err = checker.FindExisting(context.Background(), "stream-1", "idem-2")
	require.NoError(t, err)
	require.False(t, ok)
}

type mockEventStore struct {
	byIdem map[string]domain.Event
	events []domain.Event
}

func (m *mockEventStore) PutEvent(_ context.Context, e domain.Event) error {
	m.events = append(m.events, e)
	if m.byIdem == nil {
		m.byIdem = map[string]domain.Event{}
	}
	if e.IdempotencyKey != "" {
		m.byIdem[e.StreamID+"#"+e.IdempotencyKey] = e
	}
	return nil
}

func (m *mockEventStore) PutEventsBatch(context.Context, []domain.Event) error { return nil }

func (m *mockEventStore) GetByEventID(context.Context, string) (domain.Event, error) {
	return domain.Event{}, fmt.Errorf("not implemented")
}

func (m *mockEventStore) FindByIdempotencyKey(_ context.Context, streamID, key string) (domain.Event, error) {
	e, ok := m.byIdem[streamID+"#"+key]
	if !ok {
		return domain.Event{}, fmt.Errorf("missing: %w", domain.ErrNotFound)
	}
	return e, nil
}

func (m *mockEventStore) GetLatestSequence(_ context.Context, streamID string) (int64, error) {
	latest := int64(0)
	for _, event := range m.events {
		if event.StreamID == streamID && event.SequenceNumber > latest {
			latest = event.SequenceNumber
		}
	}
	return latest, nil
}

func (m *mockEventStore) QueryByStream(_ context.Context, streamID string, fromSeq int64, _ string, limit int32) ([]domain.Event, int64, bool, error) {
	if limit <= 0 {
		limit = 50
	}
	collected := make([]domain.Event, 0, limit)
	for _, event := range m.events {
		if event.StreamID == streamID && event.SequenceNumber >= fromSeq {
			collected = append(collected, event)
			if len(collected) == int(limit) {
				break
			}
		}
	}
	hasMore := len(collected) == int(limit)
	next := fromSeq + int64(len(collected))
	return collected, next, hasMore, nil
}
