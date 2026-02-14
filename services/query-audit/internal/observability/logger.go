package observability

import (
	"io"
	"log/slog"
	"os"
)

// NewLogger creates a new structured logger
func NewLogger(environment string) *slog.Logger {
	var level slog.Level
	if environment == "production" {
		level = slog.LevelInfo
	} else {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}

// NewLoggerWithWriter creates a logger with custom writer
func NewLoggerWithWriter(environment string, w io.Writer) *slog.Logger {
	var level slog.Level
	if environment == "production" {
		level = slog.LevelInfo
	} else {
		level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
	})

	return slog.New(handler)
}
