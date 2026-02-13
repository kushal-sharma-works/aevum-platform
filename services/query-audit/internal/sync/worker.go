package sync

import (
	"context"
	"log/slog"
	"math"
	"time"
)

// Worker manages background sync
type Worker struct {
	serviceName string
	syncFunc    func(context.Context, string) (string, error)
	interval    time.Duration
	maxBackoff  time.Duration
	logger      *slog.Logger
	stopChan    chan struct{}
}

// NewWorker creates a new sync worker
func NewWorker(serviceName string, syncFunc func(context.Context, string) (string, error), interval, maxBackoff time.Duration, logger *slog.Logger) *Worker {
	return &Worker{
		serviceName: serviceName,
		syncFunc:    syncFunc,
		interval:    interval,
		maxBackoff:  maxBackoff,
		logger:      logger,
		stopChan:    make(chan struct{}),
	}
}

// Start begins the background sync
func (w *Worker) Start(ctx context.Context, initialCursor string) {
	go w.run(ctx, initialCursor)
}

// Stop stops the background sync
func (w *Worker) Stop() {
	close(w.stopChan)
}

// run executes the sync loop with exponential backoff
func (w *Worker) run(ctx context.Context, cursor string) {
	backoff := w.interval
	for {
		select {
		case <-w.stopChan:
			w.logger.Info("worker stopped", slog.String("service", w.serviceName))
			return
		case <-ctx.Done():
			return
		case <-time.After(backoff):
			newCursor, err := w.syncFunc(ctx, cursor)
			if err != nil {
				w.logger.Error("sync failed", slog.String("service", w.serviceName), slog.Any("error", err))
				backoff = time.Duration(math.Min(float64(backoff)*2, float64(w.maxBackoff)))
				continue
			}
			w.logger.Info("sync successful", slog.String("service", w.serviceName))
			cursor = newCursor
			backoff = w.interval
		}
	}
}
