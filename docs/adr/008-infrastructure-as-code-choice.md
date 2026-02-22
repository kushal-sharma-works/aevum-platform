# ADR-008: Infrastructure as Code Choice

## Status

Accepted

## Date

2026-02-07

## Context

The platform requires repeatable multi-environment AWS provisioning with maintainable infrastructure logic and strong developer ergonomics.

## Decision

Use Pulumi with TypeScript for infrastructure as code.

## Consequences

### Positive

- Type-safe infrastructure definitions with familiar language tooling.
- Reusable abstractions and testable component composition.

### Negative

- Smaller ecosystem than Terraform in some domains.
- Requires Pulumi state/backend operational governance.

## Risks

- **Risk**: state drift or improper state management.
  **Mitigation**: locked state backend, CI-managed deployments, and drift detection procedures.

## Alternatives Considered

### Terraform

Pros: large ecosystem and community support.
Cons: less expressive configuration language.
Why rejected: TypeScript-based abstractions were preferred.

### AWS CDK

Pros: strong AWS-native abstractions.
Cons: tighter AWS coupling.
Why rejected: Pulumi keeps future multi-cloud optionality.

### CloudFormation

Pros: native AWS support.
Cons: verbosity and limited abstraction ergonomics.
Why rejected: lower maintainability for this platform complexity.
