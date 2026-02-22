package sync

import (
	"context"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"
)

func TestSyncStateTransitions(t *testing.T) {
	state := NewSyncState("event-timeline")
	if state.ServiceName != "event-timeline" || state.SyncStatus != "initialized" {
		t.Fatal("unexpected initial sync state")
	}

	state.UpdateCursor("cursor-1")
	if state.LastSyncedCursor != "cursor-1" || state.SyncStatus != "synced" {
		t.Fatal("expected synced state")
	}

	state.MarkFailed()
	if state.SyncStatus != "failed" {
		t.Fatal("expected failed state")
	}
}

func TestWorkerRunAndStop(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	var calls atomic.Int32

	worker := NewWorker(
		"svc",
		func(_ context.Context, cursor string) (string, error) {
			calls.Add(1)
			return cursor + "x", nil
		},
		2*time.Millisecond,
		20*time.Millisecond,
		logger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker.Start(ctx, "")
	time.Sleep(12 * time.Millisecond)
	worker.Stop()

	if calls.Load() == 0 {
		t.Fatal("expected worker to execute at least once")
	}
}
