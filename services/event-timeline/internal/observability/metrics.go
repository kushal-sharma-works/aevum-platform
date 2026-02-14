package observability

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	Registry                 *prometheus.Registry
	EventsIngestedTotal      *prometheus.CounterVec
	IngestionDurationSeconds prometheus.Histogram
	ReplayDurationSeconds    prometheus.Histogram
	ReplayEventsTotal        prometheus.Counter
	ActiveReplays            prometheus.Gauge
	HTTPRequestTotal         *prometheus.CounterVec
	HTTPRequestDuration      *prometheus.HistogramVec
}

func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	m := &Metrics{
		Registry: registry,
		EventsIngestedTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "aevum_events_ingested_total",
			Help: "Total ingested events",
		}, []string{"stream_id", "event_type", "status"}),
		IngestionDurationSeconds: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "aevum_ingestion_duration_seconds",
			Help: "Ingestion duration seconds",
		}),
		ReplayDurationSeconds: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "aevum_replay_duration_seconds",
			Help: "Replay duration seconds",
		}),
		ReplayEventsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "aevum_replay_events_total",
			Help: "Total replayed events",
		}),
		ActiveReplays: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "aevum_active_replays",
			Help: "Currently active replay operations",
		}),
		HTTPRequestTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "aevum_http_requests_total",
			Help: "HTTP requests total",
		}, []string{"method", "path", "status"}),
		HTTPRequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "aevum_http_request_duration_seconds",
			Help: "HTTP request duration seconds",
		}, []string{"method", "path"}),
	}
	registry.MustRegister(
		m.EventsIngestedTotal,
		m.IngestionDurationSeconds,
		m.ReplayDurationSeconds,
		m.ReplayEventsTotal,
		m.ActiveReplays,
		m.HTTPRequestTotal,
		m.HTTPRequestDuration,
	)
	return m
}

func (m *Metrics) RecordIngest(streamID, eventType, status string) {
	m.EventsIngestedTotal.WithLabelValues(streamID, eventType, status).Inc()
}

func (m *Metrics) ObserveIngestionDuration(seconds float64) {
	m.IngestionDurationSeconds.Observe(seconds)
}

func (m *Metrics) ObserveReplayDuration(seconds float64) {
	m.ReplayDurationSeconds.Observe(seconds)
}
