package testhelpers

import (
	"encoding/json"
	"time"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
)

func EventFixture(streamID, eventType string, occurredAt time.Time) ingest.EventInput {
	return ingest.EventInput{
		StreamID:   streamID,
		EventType:  eventType,
		Payload:    json.RawMessage(`{"fixture":true}`),
		OccurredAt: occurredAt,
	}
}
