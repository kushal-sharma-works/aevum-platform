# ADR-003: Database Selection

## Status

Accepted

## Date

2026-02-07

## Context

The platform has three distinct access patterns: immutable append-only event storage, flexible versioned rule/decision documents, and temporal full-text/correlation search.

## Decision

Use DynamoDB for events, MongoDB for rules/decisions, and Elasticsearch/OpenSearch for query and audit projections.

## Consequences

### Positive

- Each data store is optimized for its domain access pattern.
- Avoids forcing incompatible workloads into a single storage model.

### Negative

- Increases operational and platform complexity.
- Requires cross-store consistency and sync pipeline discipline.

## Risks

- **Risk**: drift between source-of-truth stores and search projection.
  **Mitigation**: cursor-based sync workers, lag metrics, replay/reindex jobs.

## Alternatives Considered

### PostgreSQL for everything

Pros: operational simplicity.
Cons: weaker fit for append-only stream scale + full-text temporal search.
Why rejected: trade-offs degrade core platform characteristics.

### Single NoSQL engine

Pros: fewer operational primitives.
Cons: compromises either querying depth or event write guarantees.
Why rejected: no single engine served all three patterns adequately.
