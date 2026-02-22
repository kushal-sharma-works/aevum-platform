# Local Development Guide

## Prerequisites

- Docker + Docker Compose
- Go 1.22+
- .NET 9 SDK
- Node.js 22+

## Initial Setup

```bash
git clone <repo-url>
cd aevum-platform
docker compose up -d --build
```

This starts the full local stack and automatically runs deterministic seeding.

`seed-data` is configured as a startup prerequisite for `query-audit` and `aevum-ui`, so users get non-empty local data before opening the UI.

## Local Workflow (Compose)

Use `docker compose` as the primary local orchestrator:

```bash
docker compose up -d --build
docker compose ps
docker compose logs -f seed-data event-timeline decision-engine query-audit
```

Open the app at `http://localhost:3000` once `seed-data` has completed.

Expected immediately after first startup:
- Rules page: non-zero active rules.
- Decisions page: non-zero seeded decisions.
- Events page: non-zero events for stream `default`.
- Timeline page: non-zero recent event entries.

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

Notes:
- `query-audit` and `aevum-ui` are started by root `docker-compose.yml`.
- Seeding is automatic on startup via the `seed-data` service.

## Viewing Traces and Metrics

- Prometheus/Grafana are available via optional devops assets and are not part of the default root compose startup.

## Seeding and Manual Flow Testing

```bash
docker compose run --rm seed-data
```

Then validate:

1. Open UI at `http://localhost:3000`.
2. Verify Decisions page shows seeded records.
3. Verify Events page shows seeded entries for stream `default`.
4. Verify Timeline page shows recent seeded events.

## Common Issues

- **Port conflict**: stop local processes using ports 8080/8081/9091/8000/27017.
- **Container startup failure**: inspect `docker compose logs --tail=200`.
- **Missing env vars**: verify `.env` or exported shell variables.
- **Test flakiness**: re-run isolated package tests and inspect dependency readiness.
