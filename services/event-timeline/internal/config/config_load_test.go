package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadSuccess(t *testing.T) {
	t.Setenv("AEVUM_JWT_SECRET", "secret")
	t.Setenv("AEVUM_GIN_PORT", "8081")
	t.Setenv("AEVUM_ECHO_PORT", "9091")
	t.Setenv("AEVUM_RATE_LIMIT_BURST", "120")
	t.Setenv("AEVUM_RATE_LIMIT_RATE", "75")
	t.Setenv("AEVUM_DYNAMODB_TABLE", "events")
	t.Setenv("AEVUM_OTEL_ENDPOINT", "otel:4317")

	cfg, err := Load()
	require.NoError(t, err)
	require.Equal(t, 8081, cfg.GinPort)
	require.Equal(t, 9091, cfg.EchoPort)
	require.Equal(t, 120, cfg.RateLimitBurst)
	require.Equal(t, float64(75), cfg.RateLimitPerSec)
	require.Equal(t, "events", cfg.DynamoTable)
}

func TestLoadValidationErrors(t *testing.T) {
	t.Run("missing jwt secret", func(t *testing.T) {
		t.Setenv("AEVUM_JWT_SECRET", "")
		_, err := Load()
		require.Error(t, err)
	})

	t.Run("invalid ports", func(t *testing.T) {
		t.Setenv("AEVUM_JWT_SECRET", "secret")
		t.Setenv("AEVUM_GIN_PORT", "0")
		_, err := Load()
		require.Error(t, err)
	})

	t.Run("invalid rate limits", func(t *testing.T) {
		t.Setenv("AEVUM_JWT_SECRET", "secret")
		t.Setenv("AEVUM_GIN_PORT", "8080")
		t.Setenv("AEVUM_ECHO_PORT", "9090")
		t.Setenv("AEVUM_RATE_LIMIT_BURST", "0")
		_, err := Load()
		require.Error(t, err)
	})

	t.Run("empty otel endpoint uses fallback", func(t *testing.T) {
		t.Setenv("AEVUM_JWT_SECRET", "secret")
		t.Setenv("AEVUM_OTEL_ENDPOINT", "")
		cfg, err := Load()
		require.NoError(t, err)
		require.NotEmpty(t, cfg.OTELEndpoint)
	})
}
