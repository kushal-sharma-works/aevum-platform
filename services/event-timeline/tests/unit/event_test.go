package unit

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

func TestNewEvent(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.NewEventInput
		wantErr bool
	}{
		{
			name: "valid",
			input: domain.NewEventInput{
				EventID:        "01HVALID",
				StreamID:       "stream-1",
				SequenceNumber: 1,
				EventType:      "created",
				Payload:        json.RawMessage(`{"ok":true}`),
				OccurredAt:     time.Now().UTC(),
				IngestedAt:     time.Now().UTC(),
				SchemaVersion:  1,
			},
		},
		{
			name: "missing stream",
			input: domain.NewEventInput{
				EventID:    "01HINVALID",
				EventType:  "created",
				Payload:    json.RawMessage(`{"ok":true}`),
				OccurredAt: time.Now().UTC(),
				IngestedAt: time.Now().UTC(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := domain.NewEvent(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.input.StreamID, event.StreamID)
		})
	}
}
