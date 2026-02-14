#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "[local-setup] Starting local infrastructure and services with Docker Compose..."
cd "$ROOT_DIR"
docker compose up -d --build

echo "[local-setup] Waiting for API health endpoints..."
for _ in {1..30}; do
  if curl -fsS http://localhost:8080/health >/dev/null 2>&1 && curl -fsS http://localhost:9091/admin/health >/dev/null 2>&1; then
    echo "[local-setup] Services are healthy."
    exit 0
  fi
  sleep 2
done

echo "[local-setup] Timeout waiting for services. Run: docker compose logs --tail=100"
exit 1
