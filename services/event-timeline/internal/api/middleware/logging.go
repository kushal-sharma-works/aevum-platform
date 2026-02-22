package middleware

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
)

func Logging(logger *slog.Logger, metrics *observability.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		metrics.HTTPRequestTotal.WithLabelValues(c.Request.Method, path, strconv.Itoa(status)).Inc()
		metrics.HTTPRequestDuration.WithLabelValues(c.Request.Method, path).Observe(time.Since(start).Seconds())
		logger.Info("http request",
			slog.String("request_id", c.GetString(RequestIDContextKey)),
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.Int("status", status),
			slog.Duration("duration", time.Since(start)),
		)
	}
}
