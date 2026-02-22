package ingest

import (
	"context"
	"errors"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
)

type IdempotencyChecker struct {
	store storage.EventStore
}

func NewIdempotencyChecker(store storage.EventStore) *IdempotencyChecker {
	return &IdempotencyChecker{store: store}
}

func (c *IdempotencyChecker) FindExisting(ctx context.Context, streamID, key string) (domain.Event, bool, error) {
	if key == "" {
		return domain.Event{}, false, nil
	}
	event, err := c.store.FindByIdempotencyKey(ctx, streamID, key)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.Event{}, false, nil
		}
		return domain.Event{}, false, err
	}
	return event, true, nil
}
