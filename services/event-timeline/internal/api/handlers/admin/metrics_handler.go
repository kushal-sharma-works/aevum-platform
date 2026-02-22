package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
)

type MetricsHandler struct {
	httpHandler http.Handler
}

func NewMetricsHandler(metrics *observability.Metrics) *MetricsHandler {
	return &MetricsHandler{httpHandler: promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})}
}

func (h *MetricsHandler) GetMetrics(c echo.Context) error {
	h.httpHandler.ServeHTTP(c.Response(), c.Request())
	return nil
}
