#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OVERLAY_PATH="$ROOT_DIR/devops/k8s/overlays/sit"

if ! command -v kubectl >/dev/null 2>&1; then
  echo "kubectl is required"
  exit 1
fi

echo "[sit-deploy] Applying SIT overlay..."
kubectl apply -k "$OVERLAY_PATH"

echo "[sit-deploy] Waiting for workloads in namespace aevum-sit..."
kubectl rollout status deployment/decision-engine -n aevum-sit --timeout=180s
kubectl rollout status deployment/event-timeline -n aevum-sit --timeout=180s
kubectl rollout status deployment/mongodb -n aevum-sit --timeout=180s
kubectl rollout status deployment/dynamodb-local -n aevum-sit --timeout=180s

echo "[sit-deploy] SIT deployment complete."
