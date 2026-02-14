package jsonutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMustMarshalAndUnmarshal(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	b := MustMarshal(payload{Name: "aevum"})
	var out payload
	err := Unmarshal(b, &out)
	require.NoError(t, err)
	require.Equal(t, "aevum", out.Name)
}

func TestMustMarshalPanicsForUnsupportedValue(t *testing.T) {
	require.Panics(t, func() {
		_ = MustMarshal(map[string]any{"bad": func() {}})
	})
}
