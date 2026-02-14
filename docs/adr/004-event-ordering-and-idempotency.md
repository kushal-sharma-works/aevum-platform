# ADR-004: Event Ordering and Idempotency

## Status

Accepted

## Date

2026-02-07

## Context

Events must be strictly ordered within each stream and ingest retries must not create duplicates or corrupt sequence guarantees.

## Decision

Assign monotonically increasing sequence numbers per stream server-side, enforce writes with DynamoDB conditional expressions, and use client idempotency keys for duplicate detection via secondary index lookup.

## Consequences

### Positive

- Preserves per-stream ordering invariants.
- Enables safe producer retries without duplication.

### Negative

- Adds read-before-write overhead for sequence assignment.
- Requires careful retry/backoff semantics under load.

## Risks

- **Risk**: hot streams causing contention on sequence assignment.
  **Mitigation**: bounded retries, adaptive backoff, and stream partitioning strategy.

## Alternatives Considered

### Client-assigned sequence numbers

Pros: lower server coordination overhead.
Cons: untrusted ordering and race conditions.
Why rejected: violates deterministic ordering guarantees.

### Timestamp-only ordering

Pros: simple implementation.
Cons: clock skew and tie resolution issues.
Why rejected: non-deterministic under distributed clocks.

### Kafka as primary event backbone

Pros: built-in partition ordering.
Cons: additional operational complexity for project scope.
Why rejected: DynamoDB Streams + existing architecture were sufficient.
