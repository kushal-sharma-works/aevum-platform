package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/api"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/clients"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/config"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/indexer"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/search"
	"github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/storage"
	syncpkg "github.com/kushal-sharma-works/aevum-platform/services/query-audit/internal/sync"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := observability.NewLogger(cfg.Environment)
	logger.Info("initializing service")

	// Create Elasticsearch client
	esClient, err := storage.NewElasticsearchClient(cfg.Elasticsearch.URLs)
	if err != nil {
		logger.Error("failed to create elasticsearch client", slog.Any("error", err))
		os.Exit(1)
	}

	// Check ES health
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := esClient.Health(ctx); err != nil {
		logger.Error("elasticsearch cluster unhealthy", slog.Any("error", err))
		cancel()
		os.Exit(1)
	}
	cancel()

	// Create indexes
	indexManager := storage.NewIndexManager(esClient.GetClient())
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	if err := indexManager.CreateIndexes(ctx); err != nil {
		logger.Error("failed to create indexes", slog.Any("error", err))
		cancel()
		os.Exit(1)
	}
	cancel()

	// Initialize clients
	eventTimelineClient := clients.NewEventTimelineClient(cfg.EventTimeline.BaseURL)
	decisionEngineClient := clients.NewDecisionEngineClient(cfg.DecisionEngine.BaseURL)

	// Create bulk indexer
	bulkIndexer := indexer.NewBulkIndexer(esClient.GetClient(), cfg.Sync.BatchSize, logger)

	// Create search engines
	searchEngine := search.NewEngine(esClient.GetClient(), logger)
	temporalQuery := search.NewTemporalQuery(esClient.GetClient(), logger)
	correlationQuery := search.NewCorrelationQuery(esClient.GetClient(), logger)
	diffEngine := search.NewDiffEngine(esClient.GetClient(), logger)
	auditBuilder := search.NewAuditBuilder(esClient.GetClient(), eventTimelineClient, decisionEngineClient, logger)

	// Create sync workers
	eventIndexer := indexer.NewEventIndexer(eventTimelineClient, bulkIndexer, logger)
	decisionIndexer := indexer.NewDecisionIndexer(decisionEngineClient, bulkIndexer, logger)

	eventWorker := syncpkg.NewWorker("event-timeline", eventIndexer.Sync, cfg.Sync.Interval, cfg.Sync.MaxBackoff, logger)
	decisionWorker := syncpkg.NewWorker("decision-engine", func(ctx context.Context, cursor string) (string, error) {
		now := time.Now()
		if err := decisionIndexer.Sync(ctx, now.Add(-1*time.Hour), now); err != nil {
			return cursor, err
		}
		return cursor, nil
	}, cfg.Sync.Interval, cfg.Sync.MaxBackoff, logger)

	// Start sync workers
	workersCtx, workersCancel := context.WithCancel(context.Background())
	eventWorker.Start(workersCtx, "")
	decisionWorker.Start(workersCtx, "")

	// Setup router
	router := api.SetupRouter(searchEngine, temporalQuery, correlationQuery, diffEngine, auditBuilder)

	// Create HTTP server
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("starting server", slog.String("addr", addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.Any("error", err))
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("shutting down service")
	workersCancel()

	eventWorker.Stop()
	decisionWorker.Stop()

	// Graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", slog.Any("error", err))
	}

	logger.Info("service stopped")
}
