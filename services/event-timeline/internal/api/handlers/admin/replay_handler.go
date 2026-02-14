package admin

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/replay"
)

type ReplayHandler struct {
	engine *replay.Engine
}

func NewReplayHandler(engine *replay.Engine) *ReplayHandler {
	return &ReplayHandler{engine: engine}
}

func (h *ReplayHandler) TriggerReplay(c echo.Context) error {
	var req domain.ReplayRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	ctx, cancel := context.WithCancel(c.Request().Context())
	defer cancel()
	eventsCh, errCh := h.engine.Replay(ctx, req)
	count := 0
	for range eventsCh {
		count++
	}
	if err := <-errCh; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"status": "completed", "events_replayed": count})
}
