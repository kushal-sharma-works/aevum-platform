# Operational Runbook

## High Event Ingestion Latency

### Check

- Event Timeline p95/p99 ingestion latency dashboards.
- Request rate versus error/timeout rates.

### Investigate

- DynamoDB throttling and consumed capacity.
- Pod CPU/memory pressure and saturation.

### Remediate

- Scale Event Timeline pods horizontally.
- Review DynamoDB table scaling/capacity and hot-key behavior.

## Decision Engine Errors

### Check

- Decision Engine error rate and latency metrics.
- Health/readiness probe status.

### Investigate

- MongoDB/DocumentDB connectivity and TLS.
- Rule payload/schema regressions.
- Upstream event payload compatibility.

### Remediate

- Inspect logs/traces for failing correlation IDs.
- Roll back problematic rule versions when needed.
- Restore DB connectivity and verify retries.

## Elasticsearch Sync Lag

### Check

- `aevum_sync_lag_seconds` and worker lag dashboards.

### Investigate

- OpenSearch cluster health and disk utilization.
- Query & Audit sync worker logs.

### Remediate

- Restart sync workers.
- Resolve index resource pressure.
- Scale cluster or worker concurrency.

## Pod CrashLoopBackOff

### Check

- `kubectl describe pod`
- `kubectl logs --previous`

### Investigate

- Missing secrets/config maps.
- Dependency connectivity failures.
- OOM kill events and memory limits.

### Remediate

- Correct missing config.
- Increase memory requests/limits where justified.
- Validate readiness/liveness paths and timings.

## Full System Recovery Sequence

Restore in order:

1. Infrastructure dependencies (DynamoDB, MongoDB, OpenSearch)
2. Event Timeline
3. Decision Engine
4. Query & Audit
5. Frontend

Verification checklist:

- Health endpoints green for each service.
- Successful smoke ingestion and decision evaluation.
- Search/audit query returns indexed records.
- Observability dashboards show stable telemetry.
