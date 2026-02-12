package replay

import (
	"context"
	"time"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
)

type Engine struct {
	eventStore storage.EventStore
	metrics    *observability.Metrics
}

func NewEngine(eventStore storage.EventStore, metrics *observability.Metrics) *Engine {
	return &Engine{eventStore: eventStore, metrics: metrics}
}

func (e *Engine) Replay(ctx context.Context, req domain.ReplayRequest) (<-chan domain.Event, <-chan error) {
	eventsCh := make(chan domain.Event, 100)
	errCh := make(chan error, 1)
	opts := NewOptions(req.From, req.To, req.EventTypes, int32(req.PageSize))
	start := time.Now()
	e.metrics.ActiveReplays.Inc()

	go func() {
		defer close(eventsCh)
		defer close(errCh)
		defer e.metrics.ActiveReplays.Dec()
		defer e.metrics.ObserveReplayDuration(time.Since(start).Seconds())

		sequence := int64(1)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			events, nextSeq, hasMore, err := e.eventStore.QueryByStream(ctx, req.StreamID, sequence, domain.DirectionForward, opts.PageSize)
			if err != nil {
				errCh <- err
				return
			}
			for _, event := range events {
				if !matchesTimeRange(event, opts.From, opts.To) || !matchesType(event, opts.EventTypes) {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case eventsCh <- event:
					e.metrics.ReplayEventsTotal.Inc()
				}
			}
			if !hasMore {
				return
			}
			sequence = nextSeq
		}
	}()

	return eventsCh, errCh
}

func matchesTimeRange(event domain.Event, from, to time.Time) bool {
	if !from.IsZero() && event.OccurredAt.Before(from) {
		return false
	}
	if !to.IsZero() && event.OccurredAt.After(to) {
		return false
	}
	return true
}

func matchesType(event domain.Event, allowed map[string]struct{}) bool {
	if len(allowed) == 0 {
		return true
	}
	_, ok := allowed[event.EventType]
	return ok
}
