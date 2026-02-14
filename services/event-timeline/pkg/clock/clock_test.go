package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMockClockNowReturnsCurrent(t *testing.T) {
	now := time.Date(2026, 2, 14, 10, 0, 0, 0, time.UTC)
	c := MockClock{Current: now}
	require.Equal(t, now, c.Now())
}

func TestRealClockNowIsUTC(t *testing.T) {
	now := (RealClock{}).Now()
	require.Equal(t, time.UTC, now.Location())
}
