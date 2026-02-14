package domain

import "time"

type ReplayRequest struct {
	StreamID    string    `json:"stream_id"`
	From        time.Time `json:"from"`
	To          time.Time `json:"to"`
	EventTypes  []string  `json:"event_types"`
	PageSize    int       `json:"page_size"`
	SpeedFactor float64   `json:"speed_factor"`
}
