#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${1:-aevum-sit}"

echo "Port-forwarding Aevum services from namespace: $NAMESPACE"

kubectl port-forward -n "$NAMESPACE" svc/event-timeline 8080:8080 &
kubectl port-forward -n "$NAMESPACE" svc/decision-engine 8081:8081 &
kubectl port-forward -n "$NAMESPACE" svc/query-audit 8082:8082 &
kubectl port-forward -n "$NAMESPACE" svc/aevum-ui 3000:80 &
kubectl port-forward -n aevum-monitoring svc/prometheus 9090:9090 &
kubectl port-forward -n aevum-monitoring svc/grafana 3001:3000 &

echo "All port-forwards active. Press Ctrl+C to stop."
wait
