# Data Model

## DynamoDB: `aevum-events` table

### Table Design

- Single-table design with PK/SK pattern.
- Immutable event writes only.
- On-demand billing mode.
- DynamoDB Streams enabled (`NEW_AND_OLD_IMAGES`) for fanout.

### Attributes

| Attribute | Type | Description |
|----------|------|-------------|
| `PK` | String | Event ID |
| `SK` | String | Event sort key (`EVENT#{streamId}#{sequence}`) |
| `GSI1PK` | String | Stream ID |
| `GSI1SK` | Number | Sequence number |
| `GSI2PK` | String | Idempotency key |
| `EventType` | String | Event category/type |
| `Payload` | Map/JSON | Event payload |
| `Metadata` | Map | Source metadata |
| `OccurredAt` | String (ISO 8601) | Business event timestamp |
| `IngestedAt` | String (ISO 8601) | Ingestion timestamp |
| `SchemaVersion` | Number | Event schema version |

### Indexes

- `GSI1` (`GSI1PK`, `GSI1SK`) for ordered stream queries.
- `GSI2` (`GSI2PK`) for idempotency lookup.

### Example Item

```json
{
  "PK": "evt_01J3ZQ8A9R",
  "SK": "EVENT#account-123#00000000000000000042",
  "GSI1PK": "account-123",
  "GSI1SK": 42,
  "GSI2PK": "idem-87ab",
  "EventType": "payment_received",
  "Payload": {"amount": 1200, "currency": "EUR"},
  "Metadata": {"source": "billing"},
  "OccurredAt": "2026-02-12T10:00:00Z",
  "IngestedAt": "2026-02-12T10:00:01Z",
  "SchemaVersion": 1
}
```

## MongoDB: `aevum` database

### Collection: `rules`

Rule definitions are immutable by version.

#### Key fields

- `RuleId` (logical identifier)
- `Version` (incrementing integer)
- `IsActive` (activation flag)
- `Conditions[]`, `Actions[]`
- `EffectiveFrom`, `EffectiveUntil`

#### Indexes

- Unique: `{ RuleId: 1, Version: -1 }`
- Filter/index support: `{ IsActive: 1 }`

#### Versioning strategy

Rule updates create new documents; historical versions are retained.

#### Example

```json
{
  "_id": "66b0ff5f84d1a2",
  "RuleId": "rule-risk-score",
  "Version": 4,
  "Name": "Risk scoring",
  "Description": "Reject high-risk applications",
  "IsActive": true,
  "Priority": 50,
  "Conditions": [
    {"Field": "score", "Operator": "GreaterThan", "Value": 80}
  ],
  "Actions": [
    {"Type": "StoreDecision", "Order": 1, "Parameters": {"result": "reject"}}
  ],
  "CreatedAt": "2026-02-10T14:20:00Z"
}
```

### Collection: `decisions`

Decisions are immutable evaluation outcomes with full trace data.

#### Key fields

- `DecisionId`
- `EventId`, `StreamId`
- `RuleId`, `RuleVersion`
- `Status`, `Output`
- `Trace[]` (step-level evaluation details)
- `EvaluatedAt`, `DeterministicHash`

#### Indexes

- `{ EventId: 1 }` for idempotency
- `{ StreamId: 1, EvaluatedAt: -1 }` for timeline views
- `{ RuleId: 1, RuleVersion: 1 }` for rule correlation

#### Example

```json
{
  "_id": "66b101de8f3af0",
  "DecisionId": "dec_01J3ZR0CFQ",
  "EventId": "evt_01J3ZQ8A9R",
  "StreamId": "account-123",
  "RuleId": "rule-risk-score",
  "RuleVersion": 4,
  "Status": "Rejected",
  "Output": {"reason": "score > 80"},
  "Trace": [
    {
      "Step": 1,
      "Field": "score",
      "Operator": "GreaterThan",
      "Expected": 80,
      "Actual": 91,
      "Matched": true
    }
  ],
  "EvaluatedAt": "2026-02-12T10:00:02Z",
  "DeterministicHash": "e3b0c44298fc1c149afbf4c8996fb924..."
}
```

## Elasticsearch / OpenSearch Indices

### `aevum-events`

Purpose: search and timeline projection for events.

Example document:

```json
{
  "event_id": "evt_01J3ZQ8A9R",
  "stream_id": "account-123",
  "event_type": "payment_received",
  "occurred_at": "2026-02-12T10:00:00Z",
  "payload": {"amount": 1200, "currency": "EUR"}
}
```

### `aevum-decisions`

Purpose: decision search, correlation, and diff operations.

Example document:

```json
{
  "decision_id": "dec_01J3ZR0CFQ",
  "event_id": "evt_01J3ZQ8A9R",
  "stream_id": "account-123",
  "rule_id": "rule-risk-score",
  "rule_version": 4,
  "status": "Rejected",
  "evaluated_at": "2026-02-12T10:00:02Z",
  "trace": [
    {"field": "score", "matched": true}
  ]
}
```

### `aevum-sync-state`

Purpose: cursor and watermark tracking for indexing workers.

Example document:

```json
{
  "worker": "decisions-sync",
  "last_processed_at": "2026-02-12T10:05:00Z",
  "cursor": "stream:account-123:seq:42"
}
```
