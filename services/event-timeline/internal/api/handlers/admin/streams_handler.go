package admin

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
)

type StreamsHandler struct {
	streamStore storage.StreamStore
}

func NewStreamsHandler(streamStore storage.StreamStore) *StreamsHandler {
	return &StreamsHandler{streamStore: streamStore}
}

func (h *StreamsHandler) ListStreams(c echo.Context) error {
	streams, err := h.streamStore.ListStreams(c.Request().Context(), 200)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"streams": streams})
}
