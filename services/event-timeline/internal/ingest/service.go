package ingest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/observability"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/storage"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/clock"
	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/pkg/identifier"
)

type Service struct {
	eventStore  storage.EventStore
	idempotency *IdempotencyChecker
	idGenerator identifier.Generator
	clock       clock.Clock
	metrics     *observability.Metrics
}

func NewService(eventStore storage.EventStore, idGenerator identifier.Generator, c clock.Clock, metrics *observability.Metrics) *Service {
	return &Service{
		eventStore:  eventStore,
		idempotency: NewIdempotencyChecker(eventStore),
		idGenerator: idGenerator,
		clock:       c,
		metrics:     metrics,
	}
}

func (s *Service) Ingest(ctx context.Context, in EventInput) (domain.Event, bool, error) {
	start := time.Now()
	if err := ValidateEventInput(in); err != nil {
		s.metrics.RecordIngest(in.StreamID, in.EventType, "invalid")
		return domain.Event{}, false, err
	}
	if existing, ok, err := s.idempotency.FindExisting(ctx, in.IdempotencyKey); err != nil {
		return domain.Event{}, false, fmt.Errorf("idempotency check: %w", err)
	} else if ok {
		s.metrics.RecordIngest(in.StreamID, in.EventType, "duplicate")
		s.metrics.ObserveIngestionDuration(time.Since(start).Seconds())
		return existing, false, nil
	}

	latest, err := s.eventStore.GetLatestSequence(ctx, in.StreamID)
	if err != nil {
		return domain.Event{}, false, fmt.Errorf("get latest sequence: %w", err)
	}

	for retries := 0; retries < 3; retries++ {
		eventID, err := s.idGenerator.New(s.clock.Now())
		if err != nil {
			return domain.Event{}, false, fmt.Errorf("generate event id: %w", err)
		}
		candidate, err := domain.NewEvent(domain.NewEventInput{
			EventID:        eventID,
			StreamID:       in.StreamID,
			SequenceNumber: latest + 1,
			EventType:      in.EventType,
			Payload:        in.Payload,
			Metadata:       in.Metadata,
			IdempotencyKey: in.IdempotencyKey,
			OccurredAt:     in.OccurredAt,
			IngestedAt:     s.clock.Now(),
			SchemaVersion:  in.SchemaVersion,
		})
		if err != nil {
			return domain.Event{}, false, fmt.Errorf("construct event: %w", err)
		}
		err = s.eventStore.PutEvent(ctx, candidate)
		if err == nil {
			s.metrics.RecordIngest(in.StreamID, in.EventType, "created")
			s.metrics.ObserveIngestionDuration(time.Since(start).Seconds())
			return candidate, true, nil
		}
		if errors.Is(err, domain.ErrSequenceConflict) {
			latest++
			continue
		}
		return domain.Event{}, false, fmt.Errorf("persist event: %w", err)
	}
	return domain.Event{}, false, fmt.Errorf("max retries reached for sequence assignment")
}

type BatchResult struct {
	Event   domain.Event `json:"event"`
	Status  string       `json:"status"`
	Error   string       `json:"error,omitempty"`
	Created bool         `json:"created"`
}

func (s *Service) BatchIngest(ctx context.Context, inputs []EventInput) []BatchResult {
	results := make([]BatchResult, 0, len(inputs))
	for _, in := range inputs {
		if err := ValidateEventInput(in); err != nil {
			return []BatchResult{{Status: "invalid", Error: err.Error(), Created: false}}
		}
	}
	for _, in := range inputs {
		event, created, err := s.Ingest(ctx, in)
		if err != nil {
			results = append(results, BatchResult{Status: "error", Error: err.Error(), Created: false})
			continue
		}
		status := "duplicate"
		if created {
			status = "created"
		}
		results = append(results, BatchResult{Event: event, Status: status, Created: created})
	}
	return results
}
