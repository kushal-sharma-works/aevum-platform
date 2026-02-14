package testhelpers

import (
	"context"
	"testing"
)

func StartDynamoDBLocal(_ context.Context, t *testing.T) string {
	t.Helper()
	t.Skip("DynamoDB Local helper is a placeholder; enable testcontainers in CI")
	return ""
}
