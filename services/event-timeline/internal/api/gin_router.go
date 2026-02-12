package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/handlers"
	mw "github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/middleware"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
)

type GinDependencies struct {
	Logger      *slog.Logger
	Metrics     *observability.Metrics
	JWTSecret   string
	RatePerSec  float64
	RateBurst   int
	Ingest      *handlers.IngestHandler
	BatchIngest *handlers.BatchIngestHandler
	Stream      *handlers.StreamHandler
	Event       *handlers.EventHandler
}

func NewGinRouter(deps GinDependencies) *gin.Engine {
	r := gin.New()
	r.Use(mw.RequestID())
	r.Use(mw.Recovery(deps.Logger))
	r.Use(mw.Logging(deps.Logger, deps.Metrics))
	r.Use(observability.GinOTelMiddleware("event-timeline-public"))
	r.Use(mw.RateLimit(deps.RatePerSec, deps.RateBurst))
	r.Use(mw.JWTAuth(deps.JWTSecret))

	v1 := r.Group("/api/v1")
	v1.POST("/events", deps.Ingest.Ingest)
	v1.POST("/events/batch", deps.BatchIngest.IngestBatch)
	v1.GET("/events/:eventId", deps.Event.GetByID)
	v1.GET("/streams/:streamId/events", deps.Stream.GetByStream)

	return r
}
