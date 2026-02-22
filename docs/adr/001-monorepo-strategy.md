# ADR-001: Monorepo Strategy

## Status

Accepted

## Date

2026-02-07

## Context

The platform contains multiple services (Go, .NET), a frontend, infrastructure as code, CI/CD workflows, and deployment manifests that must evolve together.

## Decision

Use a single monorepo with strict directory boundaries for services, frontend, and platform assets.

## Consequences

### Positive

- Enables atomic cross-service changes (API + implementation + infra + docs).
- Provides unified CI/CD, governance, and onboarding.

### Negative

- Increases repository size and CI complexity.
- Requires robust change detection and ownership conventions.

## Risks

- **Risk**: broad blast radius from poorly scoped changes.
  **Mitigation**: CODEOWNERS, path-based CI, and focused PR templates.

## Alternatives Considered

### Polyrepo

Pros: service autonomy, smaller repos.
Cons: cross-repo coordination overhead, version drift.
Why rejected: team scale and release model favor atomicity over autonomy.

### Git Submodules

Pros: composable repository boundaries.
Cons: operational friction and developer ergonomics issues.
Why rejected: complexity without proportional benefit at current scale.
