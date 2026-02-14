//go:build integration

package integration

import "testing"

func TestReplayCorrectnessWithFilters(t *testing.T) {
	t.Skip("integration test requires DynamoDB Local container")
}
