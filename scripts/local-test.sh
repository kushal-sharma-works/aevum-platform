#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_INTEGRATION=false

if [[ "${1:-}" == "--integration" ]]; then
  RUN_INTEGRATION=true
fi

cd "$ROOT_DIR"

echo "[local-test] Starting required dependencies..."
docker compose up -d mongodb dynamodb-local dynamodb-init

echo "[local-test] Running decision-engine unit tests..."
cd "$ROOT_DIR/services/decision-engine"
make test-unit

if [[ "$RUN_INTEGRATION" == true ]]; then
  echo "[local-test] Running decision-engine integration tests..."
  make test-integration
fi

echo "[local-test] Running event-timeline unit tests..."
cd "$ROOT_DIR/services/event-timeline"
make test-unit

if [[ "$RUN_INTEGRATION" == true ]]; then
  echo "[local-test] Running event-timeline integration tests (tagged)..."
  make test-integration
fi

echo "[local-test] Completed successfully."
