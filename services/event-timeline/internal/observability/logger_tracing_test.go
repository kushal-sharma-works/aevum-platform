package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLoggerForSupportedLevels(t *testing.T) {
	require.NotNil(t, NewLogger("debug"))
	require.NotNil(t, NewLogger("warn"))
	require.NotNil(t, NewLogger("error"))
	require.NotNil(t, NewLogger("info"))
}

func TestInitTracerProvider(t *testing.T) {
	tv, err := InitTracerProvider(context.Background(), "localhost:4317")
	require.NoError(t, err)
	require.NotNil(t, tv)
	require.NoError(t, tv.Shutdown(context.Background()))
}
