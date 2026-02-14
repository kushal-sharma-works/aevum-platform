//go:build integration

package integration

import "testing"

func TestIngestFlowWithDynamoDBLocal(t *testing.T) {
	t.Skip("integration test requires DynamoDB Local container")
}
