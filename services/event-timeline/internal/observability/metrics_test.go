package observability

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsRegistersCollectors(t *testing.T) {
	metrics := NewMetrics()
	require.NotNil(t, metrics.Registry)
	require.NotNil(t, metrics.EventsIngestedTotal)
	require.NotNil(t, metrics.HTTPRequestTotal)
}

func TestRecordIngestIncrementsCounter(t *testing.T) {
	metrics := NewMetrics()
	metrics.RecordIngest("stream-1", "created", "success")

	value := testutil.ToFloat64(metrics.EventsIngestedTotal.WithLabelValues("stream-1", "created", "success"))
	require.Equal(t, float64(1), value)
}
