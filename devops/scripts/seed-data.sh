#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:8080"

echo "=== Seeding Aevum with test data ==="

for i in $(seq 1 10); do
  STREAM=$([ $((i % 2)) -eq 0 ] && echo "stream-orders" || echo "stream-payments")
  curl -s -X POST "$BASE_URL/api/v1/events" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer dev-token" \
    -d "{
      \"stream_id\": \"$STREAM\",
      \"event_type\": \"test.event.v1\",
      \"payload\": {\"amount\": $((RANDOM % 1000)), \"currency\": \"EUR\", \"index\": $i},
      \"metadata\": {\"source\": \"seed-script\"},
      \"idempotency_key\": \"seed-$i\",
      \"occurred_at\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"
    }" >/dev/null
  echo "  Seeded event $i to $STREAM"
done

curl -s -X POST "http://localhost:8081/api/v1/rules" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "High Value Transaction",
    "description": "Flags transactions over 500 EUR",
    "conditions": [{"field": "payload.amount", "operator": "Gt", "value": 500}],
    "actions": [{"actionType": "Flag", "parameters": {"reason": "High value transaction detected"}}]
  }' >/dev/null

echo "  Seeded rule: High Value Transaction"
echo "=== Seeding complete ==="
