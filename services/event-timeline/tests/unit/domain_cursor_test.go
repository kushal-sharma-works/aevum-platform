package unit

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

func TestCursorEncodeDecodeRoundTrip(t *testing.T) {
	original := domain.Cursor{StreamID: "stream-1", Sequence: 42, Direction: domain.DirectionForward}
	encoded := original.Encode()

	decoded, err := domain.DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("expected decode to succeed, got %v", err)
	}

	if decoded.StreamID != original.StreamID || decoded.Sequence != original.Sequence || decoded.Direction != original.Direction {
		t.Fatalf("decoded cursor mismatch: %+v", decoded)
	}
}

func TestDecodeCursorRejectsInvalidBase64(t *testing.T) {
	_, err := domain.DecodeCursor("%%%")
	if err == nil {
		t.Fatal("expected invalid base64 error")
	}
}

func TestDecodeCursorRejectsInvalidDirection(t *testing.T) {
	raw := "stream-1:10:sideways"
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	_, err := domain.DecodeCursor(encoded)
	if err == nil {
		t.Fatal("expected invalid direction error")
	}
	if !strings.Contains(err.Error(), "invalid cursor direction") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestDecodeCursorRejectsInvalidFormat(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("stream-1:10"))
	_, err := domain.DecodeCursor(encoded)
	if err == nil {
		t.Fatal("expected invalid format error")
	}
}
