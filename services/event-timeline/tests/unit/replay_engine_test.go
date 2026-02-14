package unit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/replay"
)

func TestReplayEngineReplay(t *testing.T) {
	store := &mockEventStore{}
	now := time.Now().UTC()
	for i := 1; i <= 3; i++ {
		event, err := domain.NewEvent(domain.NewEventInput{
			EventID:        "event-id",
			StreamID:       "stream-1",
			SequenceNumber: int64(i),
			EventType:      "type-a",
			Payload:        json.RawMessage(`{"i":1}`),
			OccurredAt:     now.Add(time.Duration(i) * time.Minute),
			IngestedAt:     now,
			SchemaVersion:  1,
		})
		require.NoError(t, err)
		store.events = append(store.events, event)
	}

	engine := replay.NewEngine(store, observability.NewMetrics())
	eventsCh, errCh := engine.Replay(context.Background(), domain.ReplayRequest{
		StreamID:   "stream-1",
		From:       now,
		To:         now.Add(10 * time.Minute),
		EventTypes: []string{"type-a"},
		PageSize:   2,
	})

	count := 0
	for range eventsCh {
		count++
	}
	require.NoError(t, <-errCh)
	require.Equal(t, 3, count)
}
