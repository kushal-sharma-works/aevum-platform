package unit

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kushal-sharma-works/aevum-platform/services/event-timeline/internal/config"
)

func TestConfigLoadRequiresJWTSecret(t *testing.T) {
	t.Setenv("AEVUM_JWT_SECRET", "")
	_, err := config.Load()
	require.Error(t, err)
}

func TestConfigLoadRejectsInvalidRateLimit(t *testing.T) {
	t.Setenv("AEVUM_JWT_SECRET", "secret")
	t.Setenv("AEVUM_RATE_LIMIT_RATE", "0")
	_, err := config.Load()
	require.Error(t, err)
}

func TestConfigLoadValidDefaults(t *testing.T) {
	for _, k := range []string{
		"AEVUM_LOG_LEVEL",
		"AEVUM_GIN_PORT",
		"AEVUM_ECHO_PORT",
		"AEVUM_DYNAMODB_ENDPOINT",
		"AEVUM_DYNAMODB_TABLE",
		"AEVUM_AWS_REGION",
		"AEVUM_OTEL_ENDPOINT",
		"AEVUM_RATE_LIMIT_BURST",
		"AEVUM_RATE_LIMIT_RATE",
	} {
		_ = os.Unsetenv(k)
	}
	t.Setenv("AEVUM_JWT_SECRET", "secret")
	cfg, err := config.Load()
	require.NoError(t, err)
	require.Equal(t, 8080, cfg.GinPort)
	require.Equal(t, 9090, cfg.EchoPort)
	require.Equal(t, "aevum-events", cfg.DynamoTable)
}
