package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/httputil"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
)

type StreamHandler struct {
	eventStore storage.EventStore
}

func NewStreamHandler(eventStore storage.EventStore) *StreamHandler {
	return &StreamHandler{eventStore: eventStore}
}

func (h *StreamHandler) GetByStream(c *gin.Context) {
	streamID := c.Param("streamId")
	limit := int32(50)
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			httputil.BadRequest(c, "invalid_limit", "limit must be a valid integer")
			return
		}
		if parsed > 200 {
			parsed = 200
		}
		if parsed > 0 {
			limit = int32(parsed)
		}
	}
	direction := c.DefaultQuery("direction", domain.DirectionForward)
	if direction != domain.DirectionForward && direction != domain.DirectionBackward {
		httputil.BadRequest(c, "invalid_direction", "direction must be forward or backward")
		return
	}
	fromSeq := int64(1)
	if cursorStr := c.Query("cursor"); cursorStr != "" {
		cursor, err := domain.DecodeCursor(cursorStr)
		if err != nil {
			httputil.BadRequest(c, "invalid_cursor", err.Error())
			return
		}
		if cursor.StreamID != streamID {
			httputil.BadRequest(c, "invalid_cursor", "cursor stream_id does not match request stream_id")
			return
		}
		fromSeq = cursor.Sequence
		if fromSeq < 1 {
			fromSeq = 1
		}
		direction = cursor.Direction
	}

	events, nextSeq, hasMore, err := h.eventStore.QueryByStream(c.Request.Context(), streamID, fromSeq, direction, limit)
	if err != nil {
		httputil.Internal(c, "stream_query_failed", "failed to query stream events")
		return
	}
	nextCursor := ""
	if hasMore {
		nextCursor = domain.Cursor{StreamID: streamID, Sequence: nextSeq, Direction: direction}.Encode()
	}
	c.JSON(http.StatusOK, gin.H{
		"events":      events,
		"next_cursor": nextCursor,
		"has_more":    hasMore,
	})
}
