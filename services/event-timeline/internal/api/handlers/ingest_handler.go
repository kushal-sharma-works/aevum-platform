package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/httputil"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
)

type IngestHandler struct {
	service *ingest.Service
}

func NewIngestHandler(service *ingest.Service) *IngestHandler {
	return &IngestHandler{service: service}
}

func (h *IngestHandler) Ingest(c *gin.Context) {
	var req ingest.EventInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "invalid_request", err.Error())
		return
	}
	event, created, err := h.service.Ingest(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			httputil.BadRequest(c, "validation_failed", err.Error())
			return
		}
		slog.Error("ingest failed", slog.String("error", err.Error()))
		httputil.Internal(c, "ingest_failed", "failed to ingest event")
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	c.JSON(status, gin.H{"event": event, "created": created})
}
