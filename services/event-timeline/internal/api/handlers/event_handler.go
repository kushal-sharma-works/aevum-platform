package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/httputil"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
)

type EventHandler struct {
	eventStore storage.EventStore
}

func NewEventHandler(eventStore storage.EventStore) *EventHandler {
	return &EventHandler{eventStore: eventStore}
}

func (h *EventHandler) GetByID(c *gin.Context) {
	eventID := c.Param("eventId")
	event, err := h.eventStore.GetByEventID(c.Request.Context(), eventID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			httputil.NotFound(c, "event_not_found", "event not found")
			return
		}
		httputil.Internal(c, "event_fetch_failed", "failed to fetch event")
		return
	}
	c.JSON(200, gin.H{"event": event})
}
