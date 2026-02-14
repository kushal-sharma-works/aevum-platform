# Query & Audit Service

The Query & Audit Service is the **search, correlation, and explainability layer** of the Aevum platform. It indexes events and decisions from the Event Timeline and Decision Engine services into Elasticsearch, providing fast temporal queries, full-text search, decision diffing, and complete audit trails.

## Features

- **📊 Unified Search**: Full-text search across events and decisions
- **⏱️ Temporal Queries**: Time-range queries with configurable granularity  
- **🔗 Correlation**: Find related events/decisions by stream, rule, or event ID
- **📈 Decision Diffing**: Compare decisions at two points in time or rule versions
- **🔍 Audit Trails**: Complete causal chains from event → rule → decision
-**🔄 Background Sync**: Automatic indexing from source services

## Architecture

```
┌─────────────────┐       ┌───────────────────┐       ┌────────────────┐
│ Event Timeline  │──────▶│ Query & Audit     │◀──────│ Decision Engine│
│     Service     │       │    Service        │       │    Service     │
└─────────────────┘       └────────┬──────────┘       └────────────────┘
                                   │
                                   ▼
                           ┌──────────────┐
                           │Elasticsearch │
                           │  - aevum-events     │
                           │  - aevum-decisions  │
                           │  - aevum-sync-state │
                           └──────────────┘
```

## API Endpoints

### Search

**GET `/api/v1/search`**

General full-text search across events and decisions.

**Query Parameters**:
- `q` (required): Search query
- `type` (optional): "events", "decisions", or "all" (default: "all")
- `stream_id` (optional): Filter by stream ID
- `page` (optional): Page number (default: 1)
- `size` (optional): Page size (default: 50)

**Example**:
```bash
curl "http://localhost:8080/api/v1/search?q=payment&type=all&page=1&size=20"
```

### Temporal Queries

**GET `/api/v1/timeline`**

Query events and decisions within a time range.

**Query Parameters**:
- `from` (required): Start time (RFC3339)
- `to` (required): End time (RFC3339)
- `stream_id` (optional): Filter by stream ID
- `type` (optional): "events", "decisions", or "all" (default: "all")
- `page` (optional): Page number
- `size` (optional): Page size

**Example**:
```bash
curl "http://localhost:8080/api/v1/timeline?from=2025-01-01T00:00:00Z&to=2025-12-31T23:59:59Z&type=all"
```

### Correlation

**GET `/api/v1/correlate`**

Find related events and decisions.

**Query Parameters** (at least one required):
- `event_id`: Find all decisions triggered by an event
- `decision_id`: Find the event and rule for a decision
- `rule_id`: Find all decisions made by a rule
- `rule_version`: Filter by rule version (requires `rule_id`)
- `stream_id`: Find all events/decisions in a stream
- `event_type`: Filter by event type
- `page`, `size`: Pagination

**Examples**:
```bash
# Find all decisions from rule version 2
curl "http://localhost:8080/api/v1/correlate?rule_id=payment-rules&rule_version=2"

# Find the event that triggered a decision
curl "http://localhost:8080/api/v1/correlate?decision_id=dec_123"
```

### Decision Diff

**GET `/api/v1/diff`**

Compare decisions at two points in time or between rule versions.

**Query Parameters**:
- `t1` (required): First timestamp (RFC3339)
- `t2` (required): Second timestamp (RFC3339)
- `stream_id` (optional): Filter by stream
- `rule_id` (optional): Filter by rule
- `v1`, `v2` (optional): Compare rule versions (requires `rule_id`)

**Example**:
```bash
# Compare decisions between two timestamps
curl "http://localhost:8080/api/v1/diff?t1=2025-01-01T00:00:00Z&t2=2025-01-02T00:00:00Z&stream_id=stream_1"

# Compare decisions from two rule versions
curl "http://localhost:8080/api/v1/diff?rule_id=payment-rules&v1=1&v2=2"
```

**Response**:
```json
{
  "added": [ /* new decisions */ ],
  "removed": [ /* removed decisions */ ],
  "changed": [
    {
      "decision_id": "dec_123",
      "field": "output",
      "old_value": { "approved": false },
      "new_value": { "approved": true }
    }
  ],
  "summary": {
    "total_added": 10,
    "total_removed": 2,
    "total_changed": 5
  }
}
```

### Audit Trail

**GET `/api/v1/audit/:decisionId`**

Get the complete causal chain for a decision.

**Example**:
```bash
curl "http://localhost:8080/api/v1/audit/dec_123"
```

**Response**:
```json
{
  "decision": { /* decision details */ },
  "event": { /* triggering event */ },
  "rule_definition": { /* rule used */ },
  "chain": [
    {
      "type": "event",
      "description": "Event evt_456 of type payment.received occurred at 2025-01-15T10:00:00Z",
      "data": { /* event payload */ },
      "timestamp": "2025-01-15T10:00:00Z"
    },
    {
      "type": "rule",
      "description": "Rule version 2 triggered",
      "data": { /* rule definition */ },
      "timestamp": "2025-01-15T10:00:01Z"
    },
    {
      "type": "condition",
      "description": "amount > 100: true",
      "data": { /* trace entry */ },
      "timestamp": "2025-01-15T10:00:01Z"
    },
    {
      "type": "output",
      "description": "Decision output",
      "data": { "approved": true },
      "timestamp": "2025-01-15T10:00:01Z"
    }
  ],
  "assembled_at": "2025-01-15T10:05:00Z"
}
```

### Admin Endpoints

**GET `/admin/health`** - Health check

**POST `/admin/sync`** - Trigger manual sync

**GET `/admin/metrics`** - Prometheus metrics

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `8080` |
| `ELASTICSEARCH_URLS` | Elasticsearch URLs (comma-separated) | `http://localhost:9200` |
| `EVENT_TIMELINE_URL` | Event Timeline Service base URL | `http://localhost:8081` |
| `DECISION_ENGINE_URL` | Decision Engine Service base URL | `http://localhost:5000` |
| `SYNC_INTERVAL` | Sync interval in seconds | `5` |
| `SYNC_MAX_BACKOFF` | Max backoff on sync failure (seconds) | `300` |

## Background Sync

Two sync workers run continuously:

1. **Event Sync Worker**: Polls Event Timeline every 5 seconds for new events
2. **Decision Sync Worker**: Polls Decision Engine every 5 seconds for new decisions

Sync state (last cursor, last timestamp) is persisted in the `aevum-sync-state` Elasticsearch index.

On failure, workers use exponential backoff up to 5 minutes.

## Development

### Prerequisites

- Go 1.22+
- Elasticsearch 8.x
- Running Event Timeline Service
- Running Decision Engine Service

### Build & Run

```bash
# Download dependencies
make deps

# Build
make build

# Run locally
make run

# Or with custom config
ELASTICSEARCH_URLS=http://localhost:9200 \
EVENT_TIMELINE_URL=http://localhost:8081 \
DECISION_ENGINE_URL=http://localhost:5000 \
go run ./cmd/server
```

### Testing

```bash
# Run all tests
make test

# Unit tests only
make test-unit

# Integration tests (requires Elasticsearch)
make test-integration
```

### Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

## Elasticsearch Index Design

### aevum-events

Stores indexed events from Event Timeline.

**Mappings**:
- `event_id` (keyword): Unique event ID
- `stream_id` (keyword): Stream ID
- `sequence_number` (long): Sequence number
- `event_type` (keyword): Event type
- `payload` (object): Event payload
- `metadata` (object): Event metadata
- `occurred_at` (date): Event timestamp
- `ingested_at` (date): Ingestion timestamp
- `schema_version` (integer): Schema version

### aevum-decisions

Stores indexed decisions from Decision Engine.

**Mappings**:
- `decision_id` (keyword): Unique decision ID
- `event_id` (keyword): Triggering event ID
- `stream_id` (keyword): Stream ID
- `rule_id` (keyword): Rule ID
- `rule_version` (integer): Rule version
- `status` (keyword): Decision status
- `deterministic_hash` (keyword): Deterministic hash
- `input` (object): Decision input
- `output` (object): Decision output
- `trace` (nested): Evaluation trace
- `evaluated_at` (date): Evaluation timestamp
- `event_occurred_at` (date): Event timestamp

### aevum-sync-state

Stores sync worker state.

**Mappings**:
- `source_id` (keyword): Source identifier
- `source_type` (keyword): "event_timeline" or "decision_engine"
- `last_cursor` (long): Last synced cursor
- `last_timestamp` (date): Last synced timestamp
- `updated_at` (date): State update timestamp

## Performance

- **Search latency**: ~50-200ms depending on query complexity
- **Indexing throughput**: ~500 docs/sec per sync worker
- **Storage**: ~1KB per event, ~2KB per decision

## Monitoring

Prometheus metrics exposed at `/admin/metrics`:

- `aevum_query_duration_seconds` - Query execution time
- `aevum_index_operations_total` - Total indexing operations
- `aevum_sync_lag_seconds` - Sync lag from source services
- `aevum_search_results_total` - Total search results returned

## Troubleshooting

### High sync lag

- Check Event Timeline/Decision Engine availability
- Check Elasticsearch cluster health
- Review sync worker logs for errors

### Slow queries

- Add indexes on frequently filtered fields
- Reduce page size
- Use more specific filters
- Check Elasticsearch cluster resources

### Missing data

- Verify sync workers are running (`docker logs query-audit`)
- Check sync state: `curl http://localhost:8080/admin/sync`
- Trigger manual sync: `curl -X POST http://localhost:8080/admin/sync`

## License

MIT
