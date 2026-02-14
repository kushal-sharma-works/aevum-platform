package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetEnvFallbackAndOverride(t *testing.T) {
	t.Setenv("AEVUM_TEST_ENV", "")
	require.Equal(t, "fallback", getEnv("AEVUM_TEST_ENV", "fallback"))

	t.Setenv("AEVUM_TEST_ENV", "custom")
	require.Equal(t, "custom", getEnv("AEVUM_TEST_ENV", "fallback"))
}

func TestGetEnvIntFallbackOnInvalid(t *testing.T) {
	t.Setenv("AEVUM_TEST_INT", "invalid")
	require.Equal(t, 42, getEnvInt("AEVUM_TEST_INT", 42))

	t.Setenv("AEVUM_TEST_INT", "17")
	require.Equal(t, 17, getEnvInt("AEVUM_TEST_INT", 42))
}
