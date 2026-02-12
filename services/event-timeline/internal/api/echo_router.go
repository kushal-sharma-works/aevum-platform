package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/handlers/admin"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
)

type EchoDependencies struct {
	Health  *admin.HealthHandler
	Ready   *admin.ReadyHandler
	Replay  *admin.ReplayHandler
	Streams *admin.StreamsHandler
	Metrics *admin.MetricsHandler
}

func NewEchoRouter(deps EchoDependencies) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(observability.EchoOTelMiddleware("event-timeline-admin"))

	adminGroup := e.Group("/admin")
	adminGroup.GET("/health", deps.Health.GetHealth)
	adminGroup.GET("/ready", deps.Ready.GetReady)
	adminGroup.POST("/replay", deps.Replay.TriggerReplay)
	adminGroup.GET("/streams", deps.Streams.ListStreams)
	adminGroup.GET("/metrics", deps.Metrics.GetMetrics)

	return e
}
