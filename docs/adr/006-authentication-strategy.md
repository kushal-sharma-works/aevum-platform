# ADR-006: Authentication Strategy

## Status

Accepted

## Date

2026-02-07

## Context

Public APIs require authentication while internal/admin interfaces should remain low-friction for operations and platform automation.

## Decision

Use JWT (HS256) on public APIs and expose admin APIs on separate, internal-only network paths.

## Consequences

### Positive

- Stateless authentication with predictable runtime behavior.
- Separation of public and admin planes reduces external exposure.

### Negative

- Shared secret lifecycle management is critical.
- No first-class federated identity in current scope.

## Risks

- **Risk**: secret leakage or weak rotation discipline.
  **Mitigation**: secret manager integration, regular rotation, short token TTL.

## Alternatives Considered

### API keys

Pros: simple.
Cons: weaker claims model and rotation ergonomics.
Why rejected: JWT provides better extensibility.

### OAuth2/OIDC

Pros: enterprise identity integration.
Cons: external IdP complexity beyond current scope.
Why rejected: out of scope for current platform maturity stage.

### mTLS everywhere

Pros: strong service identity.
Cons: certificate lifecycle complexity.
Why rejected: disproportionate complexity for current needs.
