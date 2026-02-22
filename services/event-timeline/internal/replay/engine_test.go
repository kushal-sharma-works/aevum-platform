package replay

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
)

type replayStore struct {
	events []domain.Event
}

func (replayStore) PutEvent(context.Context, domain.Event) error         { return nil }
func (replayStore) PutEventsBatch(context.Context, []domain.Event) error { return nil }
func (replayStore) GetByEventID(context.Context, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (replayStore) FindByIdempotencyKey(context.Context, string, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (replayStore) GetLatestSequence(context.Context, string) (int64, error) { return 0, nil }

func (s replayStore) QueryByStream(_ context.Context, streamID string, from int64, _ string, limit int32) ([]domain.Event, int64, bool, error) {
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

func buildEvent(t *testing.T, seq int64, eventType string, occurredAt time.Time) domain.Event {
	t.Helper()
	e, err := domain.NewEvent(domain.NewEventInput{
		EventID:        "evt-replay",
		StreamID:       "stream-1",
		SequenceNumber: seq,
		EventType:      eventType,
		Payload:        json.RawMessage(`{"x":1}`),
		OccurredAt:     occurredAt,
		IngestedAt:     occurredAt,
		SchemaVersion:  1,
	})
	require.NoError(t, err)
	return e
}

func TestReplayAndCollect(t *testing.T) {
	now := time.Now().UTC()
	store := replayStore{events: []domain.Event{
		buildEvent(t, 1, "created", now),
		buildEvent(t, 2, "updated", now.Add(1*time.Minute)),
	}}

	engine := NewEngine(store, observability.NewMetrics())
	eventsCh, errCh := engine.Replay(context.Background(), domain.ReplayRequest{
		StreamID:   "stream-1",
		From:       now.Add(-1 * time.Minute),
		To:         now.Add(2 * time.Minute),
		EventTypes: []string{"created", "updated"},
		PageSize:   10,
	})

	collected := Collect(eventsCh)
	require.NoError(t, <-errCh)
	require.Len(t, collected, 2)
}

func TestReplayHelpers(t *testing.T) {
	now := time.Now().UTC()
	event := buildEvent(t, 1, "created", now)

	require.True(t, matchesTimeRange(event, now.Add(-time.Minute), now.Add(time.Minute)))
	require.False(t, matchesTimeRange(event, now.Add(time.Minute), now.Add(2*time.Minute)))

	require.True(t, matchesType(event, nil))
	require.True(t, matchesType(event, map[string]struct{}{"created": {}}))
	require.False(t, matchesType(event, map[string]struct{}{"updated": {}}))
}
