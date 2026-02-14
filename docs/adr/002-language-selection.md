# ADR-002: Language Selection

## Status

Accepted

## Date

2026-02-07

## Context

Event ingestion/replay and decision evaluation have different constraints. Ingestion favors lightweight concurrency and low overhead; decision logic favors strongly typed models and explicit deterministic boundaries.

## Decision

Use Go 1.22 (Gin + Echo) for Event Timeline and C# .NET 9 for Decision Engine.

## Consequences

### Positive

- Aligns runtime characteristics with service-specific workloads.
- Demonstrates polyglot engineering capability and deliberate boundary design.

### Negative

- Requires two backend toolchains and broader team competency.
- Slightly higher CI complexity.

## Risks

- **Risk**: inconsistent engineering standards across stacks.
  **Mitigation**: enforce linting, contract testing, and shared API spec governance.

## Alternatives Considered

### Go for both services

Pros: single backend stack.
Cons: weaker expression for rich domain modeling in decision logic.
Why rejected: .NET better fits strongly typed deterministic rule engine design goals.

### C# for both services

Pros: single backend stack.
Cons: replay streaming/concurrency ergonomics less direct than goroutine model.
Why rejected: Go is better aligned with high-throughput event/replay characteristics.

### Node.js

Pros: fast iteration.
Cons: not aligned with target stack and deterministic domain modeling goals.
Why rejected: technology choice does not match platform objectives.
