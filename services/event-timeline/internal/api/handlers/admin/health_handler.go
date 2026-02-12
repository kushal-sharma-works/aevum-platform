package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
)

type HealthHandler struct {
	eventStore storage.EventStore
}

func NewHealthHandler(eventStore storage.EventStore) *HealthHandler {
	return &HealthHandler{eventStore: eventStore}
}

func (h *HealthHandler) GetHealth(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()
	_, err := h.eventStore.GetLatestSequence(ctx, "health-probe-stream")
	status := "ok"
	if err != nil {
		status = "degraded"
	}
	return c.JSON(http.StatusOK, map[string]any{"status": status})
}
