package ingest

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/clock"
)

type testStore struct {
	events []domain.Event
	byKey  map[string]domain.Event
}

type flakyStore struct {
	testStore
	sequenceConflictsLeft   int
	returnIdempotencyOnPut  bool
	idempotencyConflictOnly bool
}

func (s *flakyStore) PutEvent(_ context.Context, e domain.Event) error {
	if s.sequenceConflictsLeft > 0 {
		s.sequenceConflictsLeft--
		return domain.ErrSequenceConflict
	}
	if s.returnIdempotencyOnPut {
		return domain.ErrIdempotencyConflict
	}
	return s.testStore.PutEvent(context.Background(), e)
}

func idempotencyMapKey(streamID, key string) string {
	return streamID + "#" + key
}

func (s *testStore) PutEvent(_ context.Context, e domain.Event) error {
	s.events = append(s.events, e)
	if s.byKey == nil {
		s.byKey = map[string]domain.Event{}
	}
	if e.IdempotencyKey != "" {
		s.byKey[idempotencyMapKey(e.StreamID, e.IdempotencyKey)] = e
	}
	return nil
}
func (s *testStore) PutEventsBatch(_ context.Context, events []domain.Event) error {
	s.events = append(s.events, events...)
	return nil
}
func (s *testStore) GetByEventID(context.Context, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (s *testStore) FindByIdempotencyKey(_ context.Context, streamID, key string) (domain.Event, error) {
	e, ok := s.byKey[idempotencyMapKey(streamID, key)]
	if !ok {
		return domain.Event{}, domain.ErrNotFound
	}
	return e, nil
}
func (s *testStore) GetLatestSequence(_ context.Context, streamID string) (int64, error) {
	var latest int64
	for _, e := range s.events {
		if e.StreamID == streamID && e.SequenceNumber > latest {
			latest = e.SequenceNumber
		}
	}
	return latest, nil
}
func (s *testStore) QueryByStream(_ context.Context, streamID string, from int64, _ string, limit int32) ([]domain.Event, int64, bool, error) {
	out := make([]domain.Event, 0)
	for _, e := range s.events {
		if e.StreamID == streamID && e.SequenceNumber >= from {
			out = append(out, e)
		}
	}
	if int32(len(out)) > limit {
		out = out[:limit]
	}
	next := from + int64(len(out))
	return out, next, false, nil
}

type testGenerator struct{}

func (testGenerator) New(time.Time) (string, error) { return "01HINTERNAL", nil }

func TestIngestAndBatchIngest(t *testing.T) {
	store := &testStore{byKey: map[string]domain.Event{}}
	service := NewService(store, testGenerator{}, clock.MockClock{Current: time.Now().UTC()}, observability.NewMetrics())

	event, created, err := service.Ingest(context.Background(), EventInput{
		StreamID:       "stream-1",
		EventType:      "created",
		Payload:        json.RawMessage(`{"ok":true}`),
		OccurredAt:     time.Now().UTC(),
		IdempotencyKey: "idem-1",
	})
	require.NoError(t, err)
	require.True(t, created)
	require.Equal(t, int64(1), event.SequenceNumber)

	results := service.BatchIngest(context.Background(), []EventInput{
		{StreamID: "stream-1", EventType: "updated", Payload: json.RawMessage(`{"x":1}`), OccurredAt: time.Now().UTC()},
		{StreamID: "", EventType: "invalid", Payload: json.RawMessage(`{"x":1}`), OccurredAt: time.Now().UTC()},
	})
	require.Len(t, results, 2)
	require.Equal(t, "created", results[0].Status)
	require.Equal(t, "invalid", results[1].Status)
	require.NotEmpty(t, results[1].Error)
}

func TestValidateEventInputRejectsMissingFields(t *testing.T) {
	err := ValidateEventInput(EventInput{})
	require.Error(t, err)
	require.True(t, errors.Is(err, domain.ErrValidation))
}

func TestIngestRetriesOnSequenceConflict(t *testing.T) {
	store := &flakyStore{testStore: testStore{byKey: map[string]domain.Event{}}, sequenceConflictsLeft: 1}
	service := NewService(store, testGenerator{}, clock.MockClock{Current: time.Now().UTC()}, observability.NewMetrics())

	event, created, err := service.Ingest(context.Background(), EventInput{
		StreamID:       "stream-1",
		EventType:      "created",
		Payload:        json.RawMessage(`{"ok":true}`),
		OccurredAt:     time.Now().UTC(),
		IdempotencyKey: "idem-retry",
	})

	require.NoError(t, err)
	require.True(t, created)
	require.Equal(t, int64(2), event.SequenceNumber)
}

func TestIngestHandlesIdempotencyConflict(t *testing.T) {
	store := &flakyStore{testStore: testStore{byKey: map[string]domain.Event{}}, returnIdempotencyOnPut: true}
	service := NewService(store, testGenerator{}, clock.MockClock{Current: time.Now().UTC()}, observability.NewMetrics())

	existing, err := domain.NewEvent(domain.NewEventInput{
		EventID:        "evt-existing",
		StreamID:       "stream-1",
		SequenceNumber: 1,
		EventType:      "created",
		Payload:        json.RawMessage(`{"ok":true}`),
		IdempotencyKey: "idem-conflict",
		OccurredAt:     time.Now().UTC(),
		IngestedAt:     time.Now().UTC(),
		SchemaVersion:  1,
	})
	require.NoError(t, err)
	store.byKey[idempotencyMapKey("stream-1", "idem-conflict")] = existing

	event, created, err := service.Ingest(context.Background(), EventInput{
		StreamID:       "stream-1",
		EventType:      "created",
		Payload:        json.RawMessage(`{"ok":true}`),
		OccurredAt:     time.Now().UTC(),
		IdempotencyKey: "idem-conflict",
	})

	require.NoError(t, err)
	require.False(t, created)
	require.Equal(t, "evt-existing", event.EventID)
}

func TestIngestFailsWhenIdempotencyConflictHasNoExisting(t *testing.T) {
	store := &flakyStore{testStore: testStore{byKey: map[string]domain.Event{}}, returnIdempotencyOnPut: true}
	service := NewService(store, testGenerator{}, clock.MockClock{Current: time.Now().UTC()}, observability.NewMetrics())

	_, _, err := service.Ingest(context.Background(), EventInput{
		StreamID:       "stream-1",
		EventType:      "created",
		Payload:        json.RawMessage(`{"ok":true}`),
		OccurredAt:     time.Now().UTC(),
		IdempotencyKey: "idem-missing",
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "idempotency conflict without existing event")
}
