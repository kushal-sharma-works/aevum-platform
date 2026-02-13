package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/middleware"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/search"
)

// SetupRouter sets up the HTTP router
func SetupRouter(searchEngine *search.Engine, temporalQuery *search.TemporalQuery, correlationQuery *search.CorrelationQuery, diffEngine *search.DiffEngine, auditBuilder *search.AuditBuilder) *gin.Engine {
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.RequestIDMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 group
	v1 := router.Group("/api/v1")

	// Search endpoints
	v1.POST("/search", NewSearchHandler(searchEngine).Handle)
	v1.POST("/timeline", NewTemporalHandler(temporalQuery).Handle)
	v1.POST("/correlate", NewCorrelationHandler(correlationQuery).Handle)
	v1.POST("/diff", NewDiffHandler(diffEngine).Handle)
	v1.GET("/audit/:decisionId", NewAuditHandler(auditBuilder).Handle)

	// Admin endpoints
	admin := router.Group("/admin")
	admin.POST("/sync", NewSyncHandler().Handle)
	admin.GET("/metrics", NewMetricsHandler().Handle)

	return router
}

// SearchHandler handles search requests
type SearchHandler struct {
	engine *search.Engine
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(engine *search.Engine) *SearchHandler {
	return &SearchHandler{engine: engine}
}

// Handle handles search requests
func (sh *SearchHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// TemporalHandler handles temporal queries
type TemporalHandler struct {
	query *search.TemporalQuery
}

// NewTemporalHandler creates a new temporal handler
func NewTemporalHandler(q *search.TemporalQuery) *TemporalHandler {
	return &TemporalHandler{query: q}
}

// Handle handles temporal requests
func (th *TemporalHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// CorrelationHandler handles correlation queries
type CorrelationHandler struct {
	query *search.CorrelationQuery
}

// NewCorrelationHandler creates a new correlation handler
func NewCorrelationHandler(q *search.CorrelationQuery) *CorrelationHandler {
	return &CorrelationHandler{query: q}
}

// Handle handles correlation requests
func (ch *CorrelationHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// DiffHandler handles diff queries
type DiffHandler struct {
	engine *search.DiffEngine
}

// NewDiffHandler creates a new diff handler
func NewDiffHandler(engine *search.DiffEngine) *DiffHandler {
	return &DiffHandler{engine: engine}
}

// Handle handles diff requests
func (dh *DiffHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// AuditHandler handles audit requests
type AuditHandler struct {
	builder *search.AuditBuilder
}

// NewAuditHandler creates a new audit handler
func NewAuditHandler(builder *search.AuditBuilder) *AuditHandler {
	return &AuditHandler{builder: builder}
}

// Handle handles audit requests
func (ah *AuditHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

// SyncHandler handles sync requests
type SyncHandler struct{}

// NewSyncHandler creates a new sync handler
func NewSyncHandler() *SyncHandler {
	return &SyncHandler{}
}

// Handle handles sync requests
func (sh *SyncHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "synced"})
}

// MetricsHandler handles metrics requests
type MetricsHandler struct{}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{}
}

// Handle handles metrics requests
func (mh *MetricsHandler) Handle(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"total_documents": 0,
		"indexes":         []string{},
	})
}
