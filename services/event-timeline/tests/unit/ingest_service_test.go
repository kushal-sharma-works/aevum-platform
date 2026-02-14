package unit

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/clock"
)

type fixedID struct{}

func (fixedID) New(time.Time) (string, error) { return "01HFIXED", nil }

func TestIngestServiceIngest(t *testing.T) {
	store := &mockEventStore{byIdem: map[string]domain.Event{}}
	svc := ingest.NewService(store, fixedID{}, clock.MockClock{Current: time.Now().UTC()}, observability.NewMetrics())

	event, created, err := svc.Ingest(context.Background(), ingest.EventInput{
		StreamID:       "stream-1",
		EventType:      "created",
		Payload:        json.RawMessage(`{"x":1}`),
		OccurredAt:     time.Now().UTC(),
		IdempotencyKey: "idem-1",
	})
	require.NoError(t, err)
	require.True(t, created)
	require.Equal(t, int64(1), event.SequenceNumber)
}
