# ADR-005: Replay Determinism

## Status

Accepted

## Date

2026-02-07

## Context

Replay must reproduce historical decisions reliably for audit and what-if analysis. Non-deterministic inputs (ambient time, mutable rules, side effects) break confidence.

## Decision

Enforce deterministic replay through explicit time injection, immutable/versioned rules, idempotent decision persistence, and hash verification for replay outputs.

## Consequences

### Positive

- Enables trustworthy audit and incident reconstruction.
- Supports controlled simulation against alternate rule versions.

### Negative

- Requires strict coding discipline against ambient dependencies.
- Adds overhead for trace and hash handling.

## Risks

- **Risk**: hidden side effects in evaluation logic.
  **Mitigation**: strict unit/integration tests and code review checklist for determinism contract.

## Alternatives Considered

### Full third-party event sourcing framework

Pros: rich built-in replay tooling.
Cons: framework lock-in and higher cognitive load.
Why rejected: custom implementation is smaller, explicit, and demonstrable.
