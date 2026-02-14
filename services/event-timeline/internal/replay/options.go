package replay

import "time"

type Options struct {
	From       time.Time
	To         time.Time
	EventTypes map[string]struct{}
	PageSize   int32
}

func NewOptions(from, to time.Time, eventTypes []string, pageSize int32) Options {
	filters := make(map[string]struct{}, len(eventTypes))
	for _, t := range eventTypes {
		filters[t] = struct{}{}
	}
	if pageSize <= 0 {
		pageSize = 100
	}
	return Options{From: from, To: to, EventTypes: filters, PageSize: pageSize}
}
