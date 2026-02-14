package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/httputil"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered", slog.Any("error", recovered), slog.String("request_id", c.GetString(RequestIDContextKey)))
		httputil.Internal(c, "internal_error", "internal server error")
	})
}
