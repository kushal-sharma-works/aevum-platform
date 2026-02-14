package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewEventValidatesRequiredFields(t *testing.T) {
	_, err := NewEvent(NewEventInput{})
	require.Error(t, err)
}

func TestNewEventBuildsEvent(t *testing.T) {
	now := time.Now().UTC()
	event, err := NewEvent(NewEventInput{
		EventID:        "evt-1",
		StreamID:       "stream-1",
		SequenceNumber: 1,
		EventType:      "created",
		Payload:        json.RawMessage(`{"ok":true}`),
		OccurredAt:     now,
		IngestedAt:     now,
		SchemaVersion:  1,
	})
	require.NoError(t, err)
	require.Equal(t, "stream-1", event.StreamID)
}
