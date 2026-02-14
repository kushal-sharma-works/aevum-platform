package middleware

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
)

func TestLoggingMiddlewareRecordsMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := observability.NewMetrics()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	r := gin.New()
	r.Use(RequestID())
	r.Use(Logging(logger, metrics))
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))

	require.Equal(t, http.StatusNoContent, rec.Code)
	count := testutil.ToFloat64(metrics.HTTPRequestTotal.WithLabelValues("GET", "/ping", "204"))
	require.Equal(t, float64(1), count)
}
