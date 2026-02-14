# Local Development Guide

## Prerequisites

- Docker + Docker Compose
- Go 1.22+
- .NET 9 SDK
- Node.js 22+
- Make

## Initial Setup

```bash
git clone <repo-url>
cd aevum-platform
make dev
```

This starts local services, dependencies, and observability stack.

## Local Workflow (`make dev`)

`make dev` orchestrates core service startup with backing dependencies so end-to-end flows can be exercised immediately.

## Running Services Individually (without Compose)

### Event Timeline

```bash
cd services/event-timeline
export AEVUM_JWT_SECRET=dev-secret
go run ./cmd/server
```

### Decision Engine

```bash
cd services/decision-engine
dotnet run --project src/Aevum.DecisionEngine.Api/Aevum.DecisionEngine.Api.csproj
```

### Query & Audit

```bash
cd services/query-audit
go run ./cmd/server
```

### Frontend

```bash
cd frontend/aevum-ui
npm install
npm run dev
```

## Testing by Service

### Event Timeline

```bash
cd services/event-timeline
make test-unit
make test-integration   # opt-in
```

### Decision Engine

```bash
cd services/decision-engine
make test-unit
make test-integration   # opt-in
```

### Frontend

```bash
cd frontend/aevum-ui
npm run test
```

## Accessing Local Databases

- DynamoDB Local: `http://localhost:8000`
- MongoDB: `mongodb://localhost:27017`
- Elasticsearch/OpenSearch: environment-specific local endpoint

## Viewing Traces and Metrics

- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3001`

Use dashboards to inspect ingestion latency, decision throughput/error rate, and sync lag.

## Seeding and Manual Flow Testing

```bash
make seed
```

Then validate:

1. Ingest events.
2. Confirm decisions are generated.
3. Query timeline/audit endpoints from frontend.

## Common Issues

- **Port conflict**: stop local processes using ports 3000/8080/8081/8082/9090/3001.
- **Container startup failure**: inspect `docker compose logs --tail=200`.
- **Missing env vars**: verify `.env` or exported shell variables.
- **Test flakiness**: re-run isolated package tests and inspect dependency readiness.
