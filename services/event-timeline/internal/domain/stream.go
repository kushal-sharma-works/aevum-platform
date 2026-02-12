package domain

type Stream struct {
	StreamID       string `json:"stream_id"`
	LatestSequence int64  `json:"latest_sequence"`
}
