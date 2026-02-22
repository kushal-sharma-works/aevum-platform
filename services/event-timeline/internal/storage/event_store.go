package storage

import (
	"context"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

type EventStore interface {
	PutEvent(ctx context.Context, event domain.Event) error
	PutEventsBatch(ctx context.Context, events []domain.Event) error
	GetByEventID(ctx context.Context, eventID string) (domain.Event, error)
	FindByIdempotencyKey(ctx context.Context, streamID, key string) (domain.Event, error)
	GetLatestSequence(ctx context.Context, streamID string) (int64, error)
	QueryByStream(ctx context.Context, streamID string, fromSequence int64, direction string, limit int32) ([]domain.Event, int64, bool, error)
}
