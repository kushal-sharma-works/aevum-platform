package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/api/handlers"
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
	if searchEngine != nil {
		searchHandler := handlers.NewSearchHandler(searchEngine)
		v1.GET("/search", searchHandler.Handle)
		v1.POST("/search", searchHandler.Handle)
	}
	if temporalQuery != nil {
		temporalHandler := handlers.NewTemporalHandler(temporalQuery)
		v1.GET("/timeline", temporalHandler.Handle)
		v1.POST("/timeline", temporalHandler.Handle)
	}
	if correlationQuery != nil {
		correlationHandler := handlers.NewCorrelationHandler(correlationQuery)
		v1.GET("/correlate", correlationHandler.Handle)
		v1.POST("/correlate", correlationHandler.Handle)
	}
	if diffEngine != nil {
		diffHandler := handlers.NewDiffHandler(diffEngine)
		v1.GET("/diff", diffHandler.Handle)
		v1.POST("/diff", diffHandler.Handle)
	}
	if auditBuilder != nil {
		v1.GET("/audit/:decisionId", handlers.NewAuditHandler(auditBuilder).Handle)
	}

	// Admin endpoints
	admin := router.Group("/admin")
	admin.POST("/sync", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "synced"})
	})
	admin.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"total_documents": 0,
			"indexes":         []string{},
		})
	})

	return router
}
