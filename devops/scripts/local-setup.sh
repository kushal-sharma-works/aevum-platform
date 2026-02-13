#!/usr/bin/env bash
set -euo pipefail

echo "=== Aevum Platform Local Setup ==="

command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed."; exit 1; }
command -v aws >/dev/null 2>&1 || { echo "AWS CLI is required but not installed."; exit 1; }

echo "Creating DynamoDB table..."
aws dynamodb create-table \
  --endpoint-url http://localhost:8000 \
  --table-name aevum-events \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=N \
    AttributeName=GSI2PK,AttributeType=S \
  --key-schema AttributeName=PK,KeyType=HASH AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes \
    '[{"IndexName":"stream-sequence-index","KeySchema":[{"AttributeName":"GSI1PK","KeyType":"HASH"},{"AttributeName":"GSI1SK","KeyType":"RANGE"}],"Projection":{"ProjectionType":"ALL"}},{"IndexName":"idempotency-index","KeySchema":[{"AttributeName":"GSI2PK","KeyType":"HASH"}],"Projection":{"ProjectionType":"KEYS_ONLY"}}]' \
  --billing-mode PAY_PER_REQUEST \
  --region eu-central-1 \
  2>/dev/null || echo "Table already exists"

echo "=== Setup complete ==="
