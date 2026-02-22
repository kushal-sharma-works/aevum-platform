package observability

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewLoggerAndNewLoggerWithWriter(t *testing.T) {
	if NewLogger("development") == nil {
		t.Fatal("expected non-nil logger")
	}

	buf := &bytes.Buffer{}
	logger := NewLoggerWithWriter("production", buf)
	if logger == nil {
		t.Fatal("expected non-nil writer logger")
	}

	logger.Info("hello")
	if !strings.Contains(buf.String(), "hello") {
		t.Fatalf("expected output to contain log message, got: %s", buf.String())
	}
}
