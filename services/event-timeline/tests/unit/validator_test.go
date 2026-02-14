package unit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
)

func TestValidateEventInput(t *testing.T) {
	err := ingest.ValidateEventInput(ingest.EventInput{
		StreamID:   "stream-1",
		EventType:  "created",
		Payload:    json.RawMessage(`{"a":1}`),
		OccurredAt: time.Now().UTC(),
	})
	require.NoError(t, err)

	err = ingest.ValidateEventInput(ingest.EventInput{})
	require.Error(t, err)
}
