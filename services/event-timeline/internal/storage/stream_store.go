package storage

import (
	"context"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
)

type StreamStore interface {
	ListStreams(ctx context.Context, limit int32) ([]domain.Stream, error)
}
