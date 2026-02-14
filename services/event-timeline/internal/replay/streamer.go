package replay

import "github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/domain"

func Collect(events <-chan domain.Event) []domain.Event {
	collected := make([]domain.Event, 0)
	for event := range events {
		collected = append(collected, event)
	}
	return collected
}
