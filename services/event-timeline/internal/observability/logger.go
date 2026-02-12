package observability

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

func NewLogger(level string) *slog.Logger {
	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	h := slog.NewJSONHandler(io.Writer(os.Stdout), &slog.HandlerOptions{Level: slogLevel})
	return slog.New(h)
}
