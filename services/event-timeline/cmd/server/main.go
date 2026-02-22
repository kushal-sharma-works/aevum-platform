package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang.org/x/sync/errgroup"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/handlers"
	adminhandlers "github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/api/handlers/admin"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/config"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/ingest"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/replay"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/clock"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/identifier"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "startup failed: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := observability.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	tp, err := observability.InitTracerProvider(ctx, cfg.OTELEndpoint)
	if err != nil {
		return fmt.Errorf("init tracer provider: %w", err)
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = tp.Shutdown(shutdownCtx)
	}()

	loadOptions := []func(*awsconfig.LoadOptions) error{
		awsconfig.WithRegion(cfg.AWSRegion),
	}
	if cfg.DynamoEndpoint != "" {
		loadOptions = append(loadOptions, awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("local", "local", "")))
		loadOptions = append(loadOptions, awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID {
				return aws.Endpoint{URL: cfg.DynamoEndpoint, SigningRegion: cfg.AWSRegion, HostnameImmutable: true}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return fmt.Errorf("load aws config: %w", err)
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	metrics := observability.NewMetrics()
	eventStore := storage.NewDynamoDBEventStore(dynamoClient, cfg.DynamoTable)
	streamStore := storage.NewDynamoDBStreamStore(dynamoClient, cfg.DynamoTable)
	ingestService := ingest.NewService(eventStore, identifier.NewULIDGenerator(), clock.RealClock{}, metrics)
	replayEngine := replay.NewEngine(eventStore, metrics)

	ingestHandler := handlers.NewIngestHandler(ingestService)
	batchIngestHandler := handlers.NewBatchIngestHandler(ingestService)
	streamHandler := handlers.NewStreamHandler(eventStore)
	eventHandler := handlers.NewEventHandler(eventStore)

	healthHandler := adminhandlers.NewHealthHandler(eventStore)
	readyHandler := adminhandlers.NewReadyHandler()
	replayHandler := adminhandlers.NewReplayHandler(replayEngine)
	streamsHandler := adminhandlers.NewStreamsHandler(streamStore)
	metricsHandler := adminhandlers.NewMetricsHandler(metrics)

	ginRouter := api.NewGinRouter(api.GinDependencies{
		Logger:      logger,
		Metrics:     metrics,
		JWTSecret:   cfg.JWTSecret,
		RatePerSec:  cfg.RateLimitPerSec,
		RateBurst:   cfg.RateLimitBurst,
		Ingest:      ingestHandler,
		BatchIngest: batchIngestHandler,
		Stream:      streamHandler,
		Event:       eventHandler,
	})
	echoRouter := api.NewEchoRouter(api.EchoDependencies{
		Health:  healthHandler,
		Ready:   readyHandler,
		Replay:  replayHandler,
		Streams: streamsHandler,
		Metrics: metricsHandler,
	})

	ginServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.GinPort), Handler: ginRouter}
	echoServer := &http.Server{Addr: fmt.Sprintf(":%d", cfg.EchoPort), Handler: echoRouter}

	g, gctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		logger.Info("starting gin server", slog.Int("port", cfg.GinPort))
		if err := ginServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("gin server: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		logger.Info("starting echo server", slog.Int("port", cfg.EchoPort))
		if err := echoServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("echo server: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigCh)
		select {
		case <-gctx.Done():
			return nil
		case sig := <-sigCh:
			logger.Info("shutdown signal received", slog.String("signal", sig.String()))
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := ginServer.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("shutdown gin server: %w", err)
			}
			if err := echoServer.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("shutdown echo server: %w", err)
			}
			if err := echoRouter.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("shutdown echo router: %w", err)
			}
			return nil
		}
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("server group exited: %w", err)
	}
	return nil
}
