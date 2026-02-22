#!/usr/bin/env bash
set -euo pipefail

EVENT_TIMELINE_URL="${EVENT_TIMELINE_URL:-http://localhost:8081}"
EVENT_TIMELINE_HEALTH_URL="${EVENT_TIMELINE_HEALTH_URL:-http://localhost:9091/admin/health}"
DECISION_ENGINE_URL="${DECISION_ENGINE_URL:-http://localhost:8080}"
EVENT_TIMELINE_JWT_SECRET="${EVENT_TIMELINE_JWT_SECRET:-dev-local-secret}"
SEED_SOURCE="${SEED_SOURCE:-seed-script-v2}"

wait_for_url() {
  local url="$1"
  local max_attempts="${2:-60}"
  local sleep_seconds="${3:-2}"

  for attempt in $(seq 1 "$max_attempts"); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep "$sleep_seconds"
  done

  echo "Timed out waiting for $url" >&2
  return 1
}

b64url() {
  openssl base64 -A | tr '+/' '-_' | tr -d '='
}

generate_jwt() {
  local secret="$1"
  local now exp header payload unsigned signature
  now=$(date +%s)
  exp=$((now + 86400 * 365))

  header=$(printf '{"alg":"HS256","typ":"JWT"}' | b64url)
  payload=$(printf '{"iss":"aevum-seed","sub":"seed-runner","iat":%s,"exp":%s}' "$now" "$exp" | b64url)
  unsigned="$header.$payload"
  signature=$(printf '%s' "$unsigned" | openssl dgst -sha256 -hmac "$secret" -binary | b64url)
  printf '%s.%s\n' "$unsigned" "$signature"
}

rule_exists() {
  local rule_name="$1"
  for status in 0 1 2 3; do
    if curl -fsS "$DECISION_ENGINE_URL/api/v1/rules?status=$status" | grep -q "\"name\":\"$rule_name\""; then
      return 0
    fi
  done
  return 1
}

create_rule_if_missing() {
  local rule_name="$1"
  local body="$2"

  if rule_exists "$rule_name"; then
    echo "  Rule exists, skipping: $rule_name"
    return 0
  fi

  local status
  status=$(curl -sS -o /tmp/aevum-seed-rule.out -w "%{http_code}" \
    --connect-timeout 5 \
    --max-time 20 \
    --retry 5 \
    --retry-delay 1 \
    --retry-connrefused \
    -X POST "$DECISION_ENGINE_URL/api/v1/rules" \
    -H "Content-Type: application/json" \
    -d "$body" || true)

  if [ -z "$status" ]; then
    echo "  Failed to create rule: $rule_name (transport error)" >&2
    return 1
  fi

  if [ "$status" = "200" ] || [ "$status" = "201" ]; then
    echo "  Created rule: $rule_name"
  else
    echo "  Failed to create rule: $rule_name (HTTP $status)" >&2
    cat /tmp/aevum-seed-rule.out >&2 || true
    return 1
  fi
}

activate_seeded_rules() {
  local ids
  ids=$(curl -fsS "$DECISION_ENGINE_URL/api/v1/rules?status=0" | grep -o '"id":"[^"]*"' | cut -d '"' -f4 | sort -u || true)

  if [ -z "$ids" ]; then
    echo "  No inactive rules to activate"
    return 0
  fi

  echo "Activating rules..."
  local rid status
  for rid in $ids; do
    status=$(curl -sS -o /tmp/aevum-seed-activate.out -w "%{http_code}" \
      --connect-timeout 5 \
      --max-time 20 \
      --retry 3 \
      --retry-delay 1 \
      --retry-connrefused \
      -X POST "$DECISION_ENGINE_URL/api/v1/rules/$rid/activate" || true)

    if [ "$status" = "200" ] || [ "$status" = "201" ]; then
      echo "  Activated rule: $rid"
    else
      echo "  Warning: activate failed for $rid (HTTP ${status:-n/a})"
    fi
  done
}

seed_decisions_for_active_rules() {
  local rules_json
  rules_json=$(curl -fsS "$DECISION_ENGINE_URL/api/v1/rules" || true)

  local rule_ids
  rule_ids=$(printf '%s' "$rules_json" | grep -o '"id":"[^"]*"' | cut -d '"' -f4 | sort -u || true)

  if [ -z "$rule_ids" ]; then
    echo "No active rules found for decision seeding"
    return 0
  fi

  echo "Seeding decisions for active rules..."

  local i=0
  local rid request_id status payload
  for rid in $rule_ids; do
    i=$((i + 1))
    request_id="seed-decision-${rid:0:8}-$(date +%s)-$i"

    payload=$(cat <<EOF
{
  "ruleId": "$rid",
  "context": {
    "amount": $((100 * i + 50)),
    "country": "US",
    "merchant": "prod_shop_$i",
    "channel": "api"
  },
  "requestId": "$request_id",
  "metadata": {
    "source": "$SEED_SOURCE",
    "scenario": "decision-seed"
  }
}
EOF
)

    status=$(curl -sS -o /tmp/aevum-seed-decision.out -w "%{http_code}" \
      --connect-timeout 5 \
      --max-time 20 \
      --retry 3 \
      --retry-delay 1 \
      --retry-connrefused \
      -X POST "$DECISION_ENGINE_URL/api/v1/decisions/evaluate" \
      -H "Content-Type: application/json" \
      -d "$payload" || true)

    if [ "$status" = "200" ] || [ "$status" = "201" ]; then
      echo "  Seeded decision for rule: $rid"
    else
      echo "  Warning: decision seed failed for rule $rid (HTTP ${status:-n/a})"
    fi
  done
}

seed_event() {
  local stream_id="$1"
  local event_type="$2"
  local occurred_at="$3"
  local idempotency_key="$4"
  local payload_json="$5"
  local metadata_json="$6"

  local body
  body=$(cat <<EOF
{
  "stream_id": "$stream_id",
  "event_type": "$event_type",
  "payload": $payload_json,
  "metadata": $metadata_json,
  "idempotency_key": "$idempotency_key",
  "occurred_at": "$occurred_at",
  "schema_version": 1
}
EOF
)

  local status
  status=$(curl -sS -o /tmp/aevum-seed-event.out -w "%{http_code}" \
    --connect-timeout 5 \
    --max-time 20 \
    --retry 5 \
    --retry-delay 1 \
    --retry-connrefused \
    -X POST "$EVENT_TIMELINE_URL/api/v1/events" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $SEED_JWT_TOKEN" \
    -d "$body" || true)

  if [ -z "$status" ]; then
    echo "  Failed event seed: $idempotency_key (transport error)" >&2
    return 1
  fi

  if [ "$status" = "200" ] || [ "$status" = "201" ] || [ "$status" = "409" ]; then
    echo "  Seeded/exists event: $idempotency_key"
  else
    echo "  Failed event seed: $idempotency_key (HTTP $status)" >&2
    cat /tmp/aevum-seed-event.out >&2 || true
    return 1
  fi
}

echo "=== Aevum deterministic seed start ==="
echo "Event Timeline: $EVENT_TIMELINE_URL"
echo "Event Timeline Health: $EVENT_TIMELINE_HEALTH_URL"
echo "Decision Engine: $DECISION_ENGINE_URL"

wait_for_url "$EVENT_TIMELINE_HEALTH_URL"
wait_for_url "$DECISION_ENGINE_URL/health"

SEED_JWT_TOKEN="$(generate_jwt "$EVENT_TIMELINE_JWT_SECRET")"

echo "Seeding rules..."
create_rule_if_missing "High Value Approval" '{
  "name": "High Value Approval",
  "description": "Approves transactions at or above 1000",
  "conditions": [{"field": "amount", "operator": 3, "value": 1000}],
  "actions": [{"type": 4, "parameters": {"result": "approved", "reason": "high-value"}, "order": 1}],
  "priority": 90,
  "createdBy": "seed-data",
  "metadata": {"scenario": "amount-threshold", "source": "seed-script-v2"}
}'

create_rule_if_missing "Low Value Manual Review" '{
  "name": "Low Value Manual Review",
  "description": "Marks low-value transactions for manual review",
  "conditions": [{"field": "amount", "operator": 4, "value": 100}],
  "actions": [{"type": 4, "parameters": {"result": "review", "reason": "low-value"}, "order": 1}],
  "priority": 70,
  "createdBy": "seed-data",
  "metadata": {"scenario": "amount-threshold", "source": "seed-script-v2"}
}'

create_rule_if_missing "Blocked Country" '{
  "name": "Blocked Country",
  "description": "Rejects transactions originating from blocked countries",
  "conditions": [{"field": "country", "operator": 10, "value": ["KP", "IR"]}],
  "actions": [{"type": 4, "parameters": {"result": "rejected", "reason": "blocked-country"}, "order": 1}],
  "priority": 100,
  "createdBy": "seed-data",
  "metadata": {"scenario": "geo-filter", "source": "seed-script-v2"}
}'

create_rule_if_missing "Merchant Prefix Risk" '{
  "name": "Merchant Prefix Risk",
  "description": "Flags test merchants by prefix",
  "conditions": [{"field": "merchant", "operator": 8, "value": "test_"}],
  "actions": [{"type": 3, "parameters": {"result": "flagged", "reason": "merchant-prefix"}, "order": 1}],
  "priority": 60,
  "createdBy": "seed-data",
  "metadata": {"scenario": "string-operator", "source": "seed-script-v2"}
}'

activate_seeded_rules
seed_decisions_for_active_rules

echo "Seeding event scenarios..."

RECENT_1=$(date -u -d '6 minutes ago' +"%Y-%m-%dT%H:%M:%SZ")
RECENT_2=$(date -u -d '4 minutes ago' +"%Y-%m-%dT%H:%M:%SZ")
RECENT_3=$(date -u -d '2 minutes ago' +"%Y-%m-%dT%H:%M:%SZ")

seed_event "default" "payment.received" "$RECENT_1" "seed-v2-default-001" '{"amount": 50, "currency": "EUR", "country": "DE", "merchant": "prod_shop_1", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"default\",\"case\":\"default-recent-1\"}"
seed_event "default" "payment.received" "$RECENT_2" "seed-v2-default-002" '{"amount": 1200, "currency": "USD", "country": "US", "merchant": "prod_shop_2", "channel": "api"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"default\",\"case\":\"default-recent-2\"}"
seed_event "default" "payment.received" "$RECENT_3" "seed-v2-default-003" '{"amount": 300, "currency": "EUR", "country": "FR", "merchant": "test_shop_1", "channel": "mobile"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"default\",\"case\":\"default-recent-3\"}"

seed_event "stream-payments-001" "payment.received" "2026-01-01T10:00:00Z" "seed-v2-payment-001" '{"amount": 50, "currency": "EUR", "country": "DE", "merchant": "store_alpha", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"low-value\",\"case\":\"approval-boundary-low\"}"
seed_event "stream-payments-001" "payment.received" "2026-01-01T10:05:00Z" "seed-v2-payment-002" '{"amount": 99, "currency": "EUR", "country": "DE", "merchant": "store_beta", "channel": "mobile"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"low-value\",\"case\":\"manual-review-upper\"}"
seed_event "stream-payments-002" "payment.received" "2026-01-01T10:10:00Z" "seed-v2-payment-003" '{"amount": 100, "currency": "EUR", "country": "FR", "merchant": "store_gamma", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"threshold\",\"case\":\"exact-100\"}"
seed_event "stream-payments-002" "payment.received" "2026-01-01T10:15:00Z" "seed-v2-payment-004" '{"amount": 999, "currency": "USD", "country": "US", "merchant": "store_delta", "channel": "api"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"threshold\",\"case\":\"just-below-1000\"}"
seed_event "stream-payments-003" "payment.received" "2026-01-01T10:20:00Z" "seed-v2-payment-005" '{"amount": 1000, "currency": "USD", "country": "US", "merchant": "store_epsilon", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"threshold\",\"case\":\"exact-1000\"}"
seed_event "stream-payments-003" "payment.received" "2026-01-01T10:25:00Z" "seed-v2-payment-006" '{"amount": 5000, "currency": "GBP", "country": "GB", "merchant": "store_zeta", "channel": "mobile"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"high-value\",\"case\":\"very-high\"}"
seed_event "stream-payments-004" "payment.received" "2026-01-01T10:30:00Z" "seed-v2-payment-007" '{"amount": 250, "currency": "EUR", "country": "IR", "merchant": "store_eta", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"geo\",\"case\":\"blocked-country-ir\"}"
seed_event "stream-payments-004" "payment.received" "2026-01-01T10:35:00Z" "seed-v2-payment-008" '{"amount": 250, "currency": "EUR", "country": "KP", "merchant": "store_theta", "channel": "api"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"geo\",\"case\":\"blocked-country-kp\"}"
seed_event "stream-payments-005" "payment.received" "2026-01-01T10:40:00Z" "seed-v2-payment-009" '{"amount": 300, "currency": "EUR", "country": "DE", "merchant": "test_shop_1", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"string\",\"case\":\"merchant-prefix-flag\"}"
seed_event "stream-payments-005" "payment.received" "2026-01-01T10:45:00Z" "seed-v2-payment-010" '{"amount": 300, "currency": "EUR", "country": "DE", "merchant": "prod_shop_1", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"string\",\"case\":\"merchant-prefix-pass\"}"

seed_event "stream-orders-001" "order.created" "2026-01-01T11:00:00Z" "seed-v2-order-001" '{"order_total": 149, "currency": "EUR", "country": "DE", "items": 2, "customer_tier": "bronze"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"orders\",\"case\":\"normal-order\"}"
seed_event "stream-orders-001" "order.updated" "2026-01-01T11:05:00Z" "seed-v2-order-002" '{"order_total": 179, "currency": "EUR", "country": "DE", "items": 3, "customer_tier": "silver"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"orders\",\"case\":\"order-update\"}"
seed_event "stream-orders-002" "order.cancelled" "2026-01-01T11:10:00Z" "seed-v2-order-003" '{"order_total": 79, "currency": "USD", "country": "US", "items": 1, "reason": "customer_request"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"orders\",\"case\":\"cancelled\"}"

seed_event "stream-replay-001" "payment.received" "2026-01-02T09:00:00Z" "seed-v2-replay-001" '{"amount": 1100, "currency": "EUR", "country": "DE", "merchant": "store_replay", "channel": "api", "request_id": "replay-1"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"replay\",\"case\":\"baseline\"}"
seed_event "stream-replay-001" "payment.received" "2026-01-02T09:01:00Z" "seed-v2-replay-002" '{"amount": 1100, "currency": "EUR", "country": "DE", "merchant": "store_replay", "channel": "api", "request_id": "replay-1"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"replay\",\"case\":\"duplicate-input-different-idem\"}"

seed_event "stream-edge-001" "payment.received" "2026-01-03T08:00:00Z" "seed-v2-edge-001" '{"amount": 0, "currency": "EUR", "country": "DE", "merchant": "store_zero", "channel": "web"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"edge\",\"case\":\"zero-amount\"}"
seed_event "stream-edge-001" "payment.received" "2026-01-03T08:05:00Z" "seed-v2-edge-002" '{"amount": -1, "currency": "EUR", "country": "DE", "merchant": "store_negative", "channel": "api"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"edge\",\"case\":\"negative-amount\"}"
seed_event "stream-edge-002" "payment.received" "2026-01-03T08:10:00Z" "seed-v2-edge-003" '{"amount": 999999, "currency": "JPY", "country": "JP", "merchant": "store_large", "channel": "mobile"}' "{\"source\":\"'$SEED_SOURCE'\",\"scenario\":\"edge\",\"case\":\"large-amount\"}"

echo "=== Aevum deterministic seed complete ==="
