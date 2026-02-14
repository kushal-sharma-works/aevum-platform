package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/httputil"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
)

type BatchIngestHandler struct {
	service *ingest.Service
}

func NewBatchIngestHandler(service *ingest.Service) *BatchIngestHandler {
	return &BatchIngestHandler{service: service}
}

func (h *BatchIngestHandler) IngestBatch(c *gin.Context) {
	var req []ingest.EventInput
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "invalid_request", err.Error())
		return
	}
	if len(req) == 0 || len(req) > 25 {
		httputil.BadRequest(c, "invalid_batch_size", "batch size must be between 1 and 25")
		return
	}
	results := h.service.BatchIngest(c.Request.Context(), req)
	c.JSON(http.StatusOK, gin.H{"results": results})
}
