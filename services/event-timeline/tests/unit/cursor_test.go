package unit

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

func TestCursorEncodeDecode(t *testing.T) {
	c := domain.Cursor{StreamID: "stream-1", Sequence: 42, Direction: domain.DirectionForward}
	encoded := c.Encode()
	decoded, err := domain.DecodeCursor(encoded)
	require.NoError(t, err)
	require.Equal(t, c, decoded)
}

func TestCursorDecodeInvalid(t *testing.T) {
	_, err := domain.DecodeCursor("not-base64")
	require.Error(t, err)
}
