package domain

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeCursorInvalidFormat(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("stream-only"))
	_, err := DecodeCursor(encoded)
	require.Error(t, err)
}

func TestDecodeCursorInvalidSequence(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("stream:not-a-number:forward"))
	_, err := DecodeCursor(encoded)
	require.Error(t, err)
}

func TestDecodeCursorInvalidDirection(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("stream:10:sideways"))
	_, err := DecodeCursor(encoded)
	require.Error(t, err)
}
