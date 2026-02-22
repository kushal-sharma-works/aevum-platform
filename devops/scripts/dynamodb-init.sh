#!/usr/bin/env sh
set -eu

export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-local}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-local}"
export AWS_DEFAULT_REGION="${AWS_DEFAULT_REGION:-eu-central-1}"

ENDPOINT="${DYNAMODB_ENDPOINT:-http://dynamodb-local:8000}"
TABLE="${DYNAMODB_TABLE:-aevum-events}"

describe_table() {
  aws dynamodb describe-table --endpoint-url "$ENDPOINT" --region "$AWS_DEFAULT_REGION" --table-name "$TABLE" >/dev/null 2>&1
}

if describe_table; then
  echo "DynamoDB table $TABLE already exists"
  exit 0
fi

aws dynamodb create-table \
  --endpoint-url "$ENDPOINT" \
  --region "$AWS_DEFAULT_REGION" \
  --table-name "$TABLE" \
  --attribute-definitions \
    AttributeName=PK,AttributeType=S \
    AttributeName=SK,AttributeType=S \
    AttributeName=GSI1PK,AttributeType=S \
    AttributeName=GSI1SK,AttributeType=N \
    AttributeName=GSI2PK,AttributeType=S \
  --key-schema \
    AttributeName=PK,KeyType=HASH \
    AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes '[{"IndexName":"GSI1","KeySchema":[{"AttributeName":"GSI1PK","KeyType":"HASH"},{"AttributeName":"GSI1SK","KeyType":"RANGE"}],"Projection":{"ProjectionType":"ALL"}},{"IndexName":"GSI2","KeySchema":[{"AttributeName":"GSI2PK","KeyType":"HASH"}],"Projection":{"ProjectionType":"ALL"}}]' \
  --billing-mode PAY_PER_REQUEST >/dev/null 2>&1 || true

attempt=1
max_attempts=30
while [ "$attempt" -le "$max_attempts" ]; do
  if describe_table; then
    echo "Created DynamoDB table $TABLE"
    exit 0
  fi
  sleep 1
  attempt=$((attempt + 1))
done

echo "Failed to verify DynamoDB table $TABLE after $max_attempts attempts" >&2
exit 1
