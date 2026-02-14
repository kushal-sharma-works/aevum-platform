# Architecture

## System Overview

Aevum is a distributed platform for deterministic, time-aware decisioning. It ingests immutable events, evaluates versioned rules, stores full decision traces, and provides temporal search and audit capabilities for operations, compliance, and engineering analysis.

The platform is intentionally designed around explicit determinism contracts. Event order is preserved per stream, rule versions are immutable, and evaluation outcomes can be replayed and verified. The architectural objective is not only high throughput in the present, but trustworthy reconstruction of the past.

At an implementation level, Aevum uses polyglot services aligned to domain requirements: Go for high-throughput event ingestion/replay orchestration, .NET for strongly modeled deterministic rule evaluation, and Go + Elasticsearch for temporal querying and audit exploration.

## C4 Model

### Level 1 — System Context

At system context level, Aevum is consumed by users (analysts, operators, developers) and receives events from external producers (platform services, partner systems, batch importers). The platform is deployed onto AWS managed infrastructure.

Reference diagram: `docs/diagrams/system-context.mmd`.

### Level 2 — Container Diagram

The Aevum platform consists of:

- Frontend (Vue) for operator workflows
- Event Timeline Service (Go) for ingestion and replay orchestration
- Decision Engine (.NET) for deterministic rule evaluation
- Query & Audit Service (Go) for temporal search and correlation
- DynamoDB, MongoDB/DocumentDB, and OpenSearch for persistence
- Supporting infra: Lambda + SQS for asynchronous fanout and OTel/Prometheus/Grafana for observability

Primary protocols:

- HTTP REST/JSON between frontend and services
- HTTP REST/JSON between internal services
- Native SDK/driver access to DynamoDB, MongoDB, and OpenSearch

Reference diagram: `docs/diagrams/container-diagram.mmd`.

## Data Flow

Primary decision flow:

1. External system sends event to Event Timeline Service (`POST /api/v1/events`).
2. Event Timeline validates payload, assigns stream sequence, and persists to DynamoDB with conditional write.
3. DynamoDB Streams emits change events.
4. Lambda consumes stream records and publishes fanout notifications to SQS.
5. Decision Engine consumes/receives event context (triggered or polling pattern).
6. Decision Engine loads relevant active rule versions from MongoDB.
7. Decision Engine evaluates rules deterministically against event context.
8. Decision output and full trace are persisted to MongoDB.
9. Query & Audit sync workers index events and decisions into Elasticsearch/OpenSearch.
10. Frontend queries Query & Audit APIs for timeline, correlation, and diff views.

Reference diagram: `docs/diagrams/data-flow.mmd`.

## Service Boundaries

### Event Timeline Service

- **Owns**: event ingestion contracts, per-stream sequence ordering, replay stream orchestration.
- **Does not own**: business decision semantics, rule lifecycle, search indexing.
- **Boundary rationale**: isolates high-throughput append-only event domain from mutable decision/search concerns.

### Decision Engine

- **Owns**: rule definitions and version lifecycle, deterministic evaluation logic, decision trace persistence.
- **Does not own**: event ingestion ordering, global search and analytics indexing.
- **Boundary rationale**: keeps deterministic domain model pure and testable with explicit time/rule inputs.

### Query & Audit Service

- **Owns**: temporal search, cross-service correlation, diff and audit retrieval APIs.
- **Does not own**: authoritative event or decision writes.
- **Boundary rationale**: read-optimized projections evolve independently from write-path latency constraints.

### Frontend

- **Owns**: operator experience, investigation workflows, replay and audit visualization.
- **Does not own**: business rules, persistence, or orchestration state.
- **Boundary rationale**: UI remains a consumer of explicit APIs and can evolve independently.

## Communication Patterns

- **Synchronous**: REST HTTP/JSON for request-response APIs between frontend and services, and internal service calls.
- **Asynchronous**: DynamoDB Streams → Lambda → SQS for event fanout and decoupled downstream processing.
- **Polling**: Query & Audit sync workers poll source services to maintain search indices and cursors.
- **Streaming to UI**: replay results are delivered via SSE for progressive rendering.

## Failure Modes

- **Event Timeline down**: ingestion unavailable; producers receive retryable errors (`503`).
- **Decision Engine down**: events still persist; decision processing backlog accumulates and drains after recovery.
- **Elasticsearch/OpenSearch down**: search/audit degraded; core ingest/evaluate paths remain operational.
- **DynamoDB throttling**: ingest latency rises; retry/backoff preserves correctness.
- **Network partition between services**: circuit breakers open; failed calls short-circuit and recover on connectivity restoration.

## Security

- JWT authentication for public APIs.
- Admin APIs exposed on separate internal plane (port + network isolation).
- TLS enforced for managed databases in production.
- DynamoDB encryption at rest (AWS-managed keys).
- OpenSearch constrained to VPC security boundaries.
- IRSA for least-privilege pod-to-AWS access.
- Secrets are externalized (Kubernetes Secrets / AWS Secrets Manager), never hardcoded in code.

## Scalability

- **Event Timeline**: stateless horizontal scale + DynamoDB on-demand scaling.
- **Decision Engine**: stateless compute scale; MongoDB/DocumentDB read scaling for rule access patterns.
- **Query & Audit**: stateless workers and independently scalable OpenSearch cluster.
- **Frontend**: static assets served from S3 + CloudFront for near-linear edge scale.
