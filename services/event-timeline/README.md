# Event Timeline Service

The Event Timeline Service is the ingestion and replay backbone for Aevum. It ingests immutable events into DynamoDB with deterministic stream ordering and supports cursor-based consumption and timestamp-based replay. The service runs two HTTP servers in one binary: Gin for public event APIs and Echo for internal admin/control APIs.

## Run locally

```bash
cd services/event-timeline
export AEVUM_JWT_SECRET=dev-secret
go run ./cmd/server
```

Public API listens on `:8080`, admin API on `:9090`.

## API Endpoints

### Public (Gin)

- `POST /api/v1/events`
- `POST /api/v1/events/batch`
- `GET /api/v1/events/:eventId`
- `GET /api/v1/streams/:streamId/events?cursor=<opaque>&limit=50&direction=forward`

Example ingest request:

```json
{
  "stream_id": "account-123",
  "event_type": "payment_received",
  "payload": {"amount": 1200, "currency": "EUR"},
  "metadata": {"source": "billing"},
  "idempotency_key": "idem-abc-1",
  "occurred_at": "2026-02-12T10:00:00Z",
  "schema_version": 1
}
```

### Admin (Echo)

- `GET /admin/health`
- `GET /admin/ready`
- `POST /admin/replay`
- `GET /admin/streams`
- `GET /admin/metrics`

## Environment variables

| Variable | Default | Required | Description |
|---|---|---|---|
| `AEVUM_LOG_LEVEL` | `info` | no | slog level |
| `AEVUM_GIN_PORT` | `8080` | no | public API port |
| `AEVUM_ECHO_PORT` | `9090` | no | admin API port |
| `AEVUM_DYNAMODB_ENDPOINT` | empty | no | custom DynamoDB endpoint (e.g. local) |
| `AEVUM_DYNAMODB_TABLE` | `aevum-events` | no | DynamoDB table name |
| `AEVUM_AWS_REGION` | `eu-central-1` | no | AWS region |
| `AEVUM_JWT_SECRET` | - | yes | HS256 secret for public API auth |
| `AEVUM_OTEL_ENDPOINT` | `localhost:4317` | no | OTLP gRPC endpoint |
| `AEVUM_RATE_LIMIT_BURST` | `100` | no | token bucket burst |
| `AEVUM_RATE_LIMIT_RATE` | `50` | no | token bucket sustained req/s |

## Tests

```bash
cd services/event-timeline
go test ./... -race
```

## Architecture decisions

- **Gin + Echo**: Gin is used for low-overhead hot-path ingestion APIs; Echo is used for internal admin APIs with clean grouped routing.
- **DynamoDB**: Single-table model with GSIs supports immutable event storage, stream ordering, and idempotency lookups.
- **ULID**: Time-sortable event identifiers preserve lexicographic order and improve replay/query characteristics.
