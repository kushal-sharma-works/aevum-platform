package api

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/handlers"
	adminhandlers "github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/handlers/admin"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/replay"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/clock"
)

type routerEventStore struct{}

func (routerEventStore) PutEvent(context.Context, domain.Event) error         { return nil }
func (routerEventStore) PutEventsBatch(context.Context, []domain.Event) error { return nil }
func (routerEventStore) GetByEventID(context.Context, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (routerEventStore) FindByIdempotencyKey(context.Context, string) (domain.Event, error) {
	return domain.Event{}, domain.ErrNotFound
}
func (routerEventStore) GetLatestSequence(context.Context, string) (int64, error) { return 0, nil }
func (routerEventStore) QueryByStream(context.Context, string, int64, string, int32) ([]domain.Event, int64, bool, error) {
	return nil, 0, false, nil
}

type routerStreamStore struct{}

func (routerStreamStore) ListStreams(context.Context, int32) ([]domain.Stream, error) {
	return nil, nil
}

type fixedRouterGenerator struct{}

func (fixedRouterGenerator) New(_ time.Time) (string, error) { return "01HROUTER", nil }

func TestNewGinRouterRegistersRoutes(t *testing.T) {
	metrics := observability.NewMetrics()
	store := routerEventStore{}
	service := ingest.NewService(store, fixedRouterGenerator{}, clock.MockClock{Current: time.Now().UTC()}, metrics)

	router := NewGinRouter(GinDependencies{
		Logger:      slog.Default(),
		Metrics:     metrics,
		JWTSecret:   "secret",
		RatePerSec:  10,
		RateBurst:   10,
		Ingest:      handlers.NewIngestHandler(service),
		BatchIngest: handlers.NewBatchIngestHandler(service),
		Stream:      handlers.NewStreamHandler(store),
		Event:       handlers.NewEventHandler(store),
	})

	routes := router.Routes()
	require.NotEmpty(t, routes)
}

func TestNewEchoRouterRegistersAdminRoutes(t *testing.T) {
	metrics := observability.NewMetrics()
	store := routerEventStore{}
	engine := replay.NewEngine(store, metrics)

	router := NewEchoRouter(EchoDependencies{
		Health:  adminhandlers.NewHealthHandler(store),
		Ready:   adminhandlers.NewReadyHandler(),
		Replay:  adminhandlers.NewReplayHandler(engine),
		Streams: adminhandlers.NewStreamsHandler(routerStreamStore{}),
		Metrics: adminhandlers.NewMetricsHandler(metrics),
	})

	routes := router.Routes()
	require.NotEmpty(t, routes)
}
