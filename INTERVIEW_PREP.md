# Aevum Platform — Interview Preparation Guide

> **How to use this guide:** Every question below is followed by a short, memorable answer and a deeper elaboration. Under pressure, give the short answer; if probed, expand with the elaboration. Think of the short answer as your headline and the elaboration as your evidence.

---

## Table of Contents

1. [Non-Technical Questions](#1-non-technical-questions)
2. [Architecture & Design](#2-architecture--design)
3. [Event Timeline Service (Go)](#3-event-timeline-service-go)
4. [Decision Engine (.NET)](#4-decision-engine-net)
5. [Query & Audit Service (Go)](#5-query--audit-service-go)
6. [Frontend (Vue 3 + TypeScript)](#6-frontend-vue-3--typescript)
7. [DevOps, Infrastructure & CI/CD](#7-devops-infrastructure--cicd)
8. [Observability](#8-observability)
9. [Security](#9-security)
10. [Testing](#10-testing)
11. [Data Model Deep-Dives](#11-data-model-deep-dives)
12. [What Would You Do Differently?](#12-what-would-you-do-differently)
13. [Additional Deep-Dive Questions](#13-additional-deep-dive-questions)
14. [Behavioral & Process Questions](#14-behavioral--process-questions)
15. [Worst-Case Rapid-Fire Questions](#15-worst-case-rapid-fire-questions)

---

## 1. Non-Technical Questions

---

### Q: Why is the project called "Aevum"?

**Short answer:** "Aevum" is the Latin word for an age or epoch—timeless, continuous duration. The name captures the platform's core purpose: understanding decisions across time, replaying history, and reasoning about what happened at any point in the past.

**Elaboration:** Most platforms live only in the present. "Aevum" signals that time itself is a first-class citizen here. Every decision is timestamped, every rule is versioned, every event is immutable—so you can always travel back to any epoch and reconstruct exactly what happened and why.

---

### Q: What problem does this platform actually solve?

**Short answer:** It answers the question most systems can't: *"What was the decision at time T, why was it made, and what would have happened under different rules?"*

**Elaboration:** Consider a bank that rejects a loan application. Three months later, the rules change and the same applicant would have been approved. Was the original rejection correct? With Aevum, you can replay the original event through the original rule version and see the full evaluation trace—field by field—and you can also simulate the new rules against the historical data before deploying them.

---

### Q: Who is the target user?

**Short answer:** Operations teams, compliance officers, and engineers who need auditability, incident reconstruction, and safe rule simulation in high-stakes domains such as finance, fraud detection, or policy enforcement.

**Elaboration:** The UI provides a timeline viewer for operators, a decision inspector for compliance, and a replay console for engineers. The APIs are also directly accessible for programmatic audit pipelines.

---

### Q: Why build this instead of using an existing rules engine?

**Short answer:** Existing engines (Drools, Easy Rules) evaluate rules but don't provide immutable event storage, deterministic replay, or temporal audit trails as integrated platform concepts.

**Elaboration:** A rules engine answers "what is the decision right now." Aevum answers "what was the decision at any past moment and can I prove it." The platform combines event sourcing, versioned rule evaluation, and an audit query layer in a coherent, deployable whole. That combination doesn't exist off the shelf.

---

### Q: Why a monorepo?

**Short answer:** Atomic changes. When you add a new API field, you can update the backend, the OpenAPI spec, the frontend client, the Helm chart, and the documentation in one pull request and one CI run.

**Elaboration:** See ADR-001. The alternatives—polyrepo and git submodules—were rejected because cross-repo coordination creates version drift and release friction. At this platform's scale (three services, shared contracts, shared deployment manifests), a monorepo reduces coordination costs without significant downside.

---

### Q: Why polyglot (Go + C#) instead of one language?

**Short answer:** Different problems need different tools. Go's goroutine model is ideal for streaming thousands of events; C#'s type system and .NET's TimeProvider are ideal for enforcing deterministic evaluation contracts.

**Elaboration:** See ADR-002. Event ingestion and replay are I/O-heavy, concurrent, and benefit from Go's lightweight goroutine model and direct channel-based streaming. Decision logic is CPU-bound, domain-rich, and benefits from C#'s records, immutability, and explicit dependency injection for `TimeProvider`. Forcing one language into both roles would mean compromising at least one workload.

---

### Q: Why three databases (DynamoDB, MongoDB, Elasticsearch)?

**Short answer:** Each store is optimized for one access pattern. DynamoDB for immutable append-only ordered writes, MongoDB for flexible versioned documents, Elasticsearch for temporal full-text and correlation queries.

**Elaboration:** See ADR-003. No single database serves all three patterns well. A single relational DB would struggle with Elasticsearch's temporal query semantics; a single NoSQL DB would compromise either write ordering or search depth. The trade-off is operational complexity, which is mitigated by cursor-based sync workers and re-index jobs.

---

### Q: Why Pulumi over Terraform?

**Short answer:** TypeScript-native abstractions, testable components, and type safety that catches infrastructure mistakes at compile time instead of apply time.

**Elaboration:** See ADR-008. Terraform's HCL is declarative but limited as a language—loops, conditionals, and abstractions are awkward. Pulumi TypeScript lets you write real functions (`createVpc`, `createEksCluster`) with interfaces, unit tests, and IDE support. The trade-off is a smaller community ecosystem in some domains, but the gain in maintainability was deemed worth it.

---

### Q: Why ArgoCD for deployments?

**Short answer:** GitOps: the Git repository is the single source of truth for cluster state. ArgoCD continuously reconciles what is in Git with what is running.

**Elaboration:** Manual `kubectl apply` is error-prone and leaves no audit trail. ArgoCD diffs the desired state (Helm charts in Git) against the live cluster and self-heals on drift. This means a bad deployment can be rolled back by reverting a commit, and every deployment is traceable to a commit hash.

---

### Q: What is the project's current maturity / what is production-ready vs. prototype?

**Short answer:** The core decision flow, event ingestion, replay engine, and audit APIs are fully functional locally. The Pulumi infrastructure is a complete but optional cloud path. The frontend covers all key workflows.

**Elaboration:** This is a fully designed platform with real code—not a toy. Tests exist across all services. The observability stack (OTel, Prometheus, Grafana) is wired up. What isn't production-hardened in the current scope: gRPC inter-service contracts, a schema registry, Playwright E2E tests, and chaos testing. Those are explicitly called out in the "What I Would Do Differently" section.

---

## 2. Architecture & Design

---

### Q: Walk me through the high-level architecture.

**Short answer:** Browser → Vue Frontend → three backend services → three databases, all observable via OTel/Prometheus/Grafana.

**Elaboration:**
1. The **Frontend** (Vue 3) calls three services via REST.
2. **Event Timeline** (Go, port 8081/9091) ingests events into DynamoDB, assigns sequence numbers, and streams replay results via SSE.
3. **Decision Engine** (.NET 9, port 8080) holds versioned rules in MongoDB and evaluates them deterministically.
4. **Query & Audit** (Go, port 8082) indexes events and decisions into Elasticsearch and exposes temporal search, correlation, and diff APIs.
5. All services emit traces and metrics through an OpenTelemetry Collector to Prometheus, visualized in Grafana.
6. In production, DynamoDB Streams → Lambda → SQS fans out events asynchronously to the Decision Engine.

---

### Q: Why does the Event Timeline expose two ports (8081 and 9091)?

**Short answer:** Separation of public and admin planes. Port 8081 is the public data API (JWT-protected). Port 9091 is the admin/replay API (internal network only, no public exposure).

**Elaboration:** The public API is what external producers use to ingest events. The admin API is used internally for replay triggers and operational commands. By isolating them on different ports and network boundaries, you reduce the blast radius of a compromised public token—it cannot access replay or admin operations.

---

### Q: How does data flow from an event to a decision?

**Short answer:** Producer POSTs event → Event Timeline assigns sequence and writes to DynamoDB → Decision Engine evaluates matching rules → Decision stored in MongoDB → Query & Audit sync worker indexes both into Elasticsearch.

**Elaboration:**
1. `POST /api/v1/events` reaches the Event Timeline ingest handler.
2. Idempotency check via GSI2 lookup; if duplicate, return existing event.
3. Server-side sequence number assigned; DynamoDB transactional write enforces uniqueness.
4. DynamoDB Streams → Lambda → SQS fans out (production path).
5. Decision Engine receives event context, loads the active rule from MongoDB, evaluates with `DeterministicEvaluator`, stores full trace in MongoDB.
6. Query & Audit sync workers poll Event Timeline and Decision Engine on a cursor, bulk-index into Elasticsearch indices `aevum-events` and `aevum-decisions`.

---

### Q: What happens if the Decision Engine is down?

**Short answer:** Events still ingest successfully. A decision backlog accumulates in the queue (SQS DLQ in production) and drains automatically once the service recovers.

**Elaboration:** The ingestion path (Event Timeline → DynamoDB) is decoupled from evaluation. Decisions are not required for events to be stored. The event is durable; the evaluation can happen later. The replay flow can also recover decisions retroactively if needed.

---

### Q: How do services communicate internally?

**Short answer:** HTTP REST/JSON synchronously between services; DynamoDB Streams → Lambda → SQS asynchronously for event fanout; SSE for streaming replay results to the UI.

**Elaboration:** There is no message broker (Kafka/RabbitMQ) in the current scope. Internal calls are straightforward REST with circuit-breaker-friendly retry semantics. The async fanout is handled by a Lambda trigger on the DynamoDB Stream, which puts messages on an SQS queue for the Decision Engine to consume. SSE (Server-Sent Events) was chosen for replay output because it is simple, browser-native, and avoids the complexity of WebSockets for a unidirectional stream.

---

### Q: What is the C4 model and does Aevum follow it?

**Short answer:** Yes. C4 has four levels: System Context, Containers, Components, Code. Aevum has documented diagrams at Level 1 (system-context.mmd) and Level 2 (container-diagram.mmd), and the codebase structure mirrors Level 3 (component boundaries per service).

---

## 3. Event Timeline Service (Go)

---

### Q: What does the Event Timeline service do?

**Short answer:** It is the immutable append-only event store. It accepts events, enforces per-stream ordering via sequence numbers, checks idempotency, and provides a replay engine that streams historical events.

---

### Q: What is the `Event` struct and what are all those DynamoDB tags?

**Short answer:** It is the core domain entity. The DynamoDB tags map Go fields to DynamoDB attribute names following a single-table design.

**Elaboration:**
- `PK` = Event ID (primary partition key for direct lookups).
- `SK` = `EVENT#{streamId}#{paddedSequence}` (sort key for ordering within a partition).
- `GSI1PK` / `GSI1SK` = Stream ID / Sequence number → ordered stream queries.
- `GSI2PK` = `{streamId}#{idempotencyKey}` → idempotency lookup.
- Fields like `json:"-"` hide internal DynamoDB keys from the API response.

---

### Q: How does sequence assignment work? What prevents races?

**Short answer:** Server reads the latest sequence, attempts a transactional write that conditionally fails if the sequence is already taken, and retries up to 3 times if there is a conflict.

**Elaboration (code walkthrough):**
```go
// ingest/service.go
latest, _ = s.eventStore.GetLatestSequence(ctx, in.StreamID)
for retries := 0; retries < 3; retries++ {
    candidate.SequenceNumber = latest + 1
    err = s.eventStore.PutEvent(ctx, candidate)
    if errors.Is(err, domain.ErrSequenceConflict) {
        latest++; continue
    }
    ...
}
```
In `PutEvent`, a `TransactWriteItems` call writes both the event and a sequence guard row with `attribute_not_exists` conditions. If another writer takes the same sequence concurrently, the transaction is cancelled, the error is mapped to `ErrSequenceConflict`, and the service increments and retries. Up to 3 retries handle transient hot-stream contention.

---

### Q: What is idempotency and how is it implemented?

**Short answer:** If a producer retries the same event (same stream + idempotency key), the existing event is returned instead of creating a duplicate.

**Elaboration:** The ingest path checks `IdempotencyChecker.FindExisting` via a GSI2 lookup before attempting a write. If a key already exists, the event is returned immediately with `created=false`. The transactional write also includes a conditional idempotency lock item—so even under race conditions between two concurrent identical requests, at most one succeeds and the other resolves to the existing event.

---

### Q: What is the `Clock` interface and why does it exist?

**Short answer:** To eliminate ambient time from the ingest path, making behavior deterministic and testable.

**Elaboration (code):**
```go
// pkg/clock/clock.go
type Clock interface { Now() time.Time }
type RealClock struct{}
func (RealClock) Now() time.Time { return time.Now().UTC() }
type MockClock struct{ Current time.Time }
func (m MockClock) Now() time.Time { return m.Current }
```
In production, `RealClock` is injected. In tests, `MockClock` is injected with a fixed time. This means a unit test can assert exact `IngestedAt` values without depending on wall-clock time. This is the Go equivalent of .NET's `TimeProvider`.

---

### Q: How does the replay engine work?

**Short answer:** A goroutine pages through DynamoDB events in sequence order via Go channels, filtering by time range and event type, and sends each matching event to an output channel consumed by the SSE streamer.

**Elaboration (code):**
```go
// replay/engine.go
go func() {
    for {
        events, nextSeq, hasMore, _ = e.eventStore.QueryByStream(ctx, req.StreamID, sequence, DirectionForward, pageSize)
        for _, event := range events {
            if !matchesTimeRange(event, opts.From, opts.To) { continue }
            eventsCh <- event
        }
        if !hasMore { return }
        sequence = nextSeq
    }
}()
```
The caller receives a `<-chan domain.Event`. This is idiomatic Go: the engine produces, the streamer consumes. Cancelling `ctx` stops the goroutine cleanly. Active replays are tracked in a Prometheus gauge (`metrics.ActiveReplays`).

---

### Q: Why does Gin and Echo both appear in the Event Timeline service?

**Short answer:** Gin handles the public REST API (port 8081); Echo handles the admin/internal API (port 9091). Both run as separate HTTP servers on separate ports, serving different audiences and route sets.

**Elaboration:** Separating them by framework gives a clear boundary between external and internal routes. Gin is more common for public REST APIs with strong middleware ecosystems; Echo was used for the admin API for its clean handler model. In production these can be isolated by network security groups.

---

### Q: What does `SK: "EVENT#{streamId}#{sequence:020d}"` mean?

**Short answer:** It is a zero-padded sort key that allows DynamoDB to return events in sequence order without a separate sort step.

**Elaboration:** DynamoDB sorts lexicographically. Without zero-padding, sequence 10 would sort before sequence 9 (`"10" < "9"` lexicographically). Padding to 20 digits (`%020d`) ensures numeric sort order is preserved in string comparison. The 20-digit pad supports sequences up to 99,999,999,999,999,999,999—effectively unlimited.

---

## 4. Decision Engine (.NET)

---

### Q: What does the Decision Engine do?

**Short answer:** It is the rule evaluation core. It stores versioned rules in MongoDB, evaluates them deterministically against event context, and stores a complete trace of every decision.

---

### Q: What is the `DeterministicEvaluator` and how does determinism work?

**Short answer:** It evaluates rule conditions against a context dictionary, produces a SHA-256 hash of the inputs, and never touches ambient state (no `DateTime.Now`, no random, no external calls).

**Elaboration (code):**
```csharp
// DeterministicEvaluator.cs
public EvaluationResult Evaluate(Rule rule, EvaluationContext context) {
    var isMatch = EvaluateConditions(rule.Conditions, context.Data, ...);
    var deterministicHash = ComputeHash(rule, context);
    ...
}

public string ComputeHash(Rule rule, EvaluationContext context) {
    var hashInput = new { RuleId, RuleVersion, RequestId, Context = SortDictionary(context.Data) };
    var json = JsonSerializer.Serialize(hashInput, camelCase);
    return Convert.ToHexString(SHA256.HashData(Encoding.UTF8.GetBytes(json)));
}
```
The hash is computed over `RuleId + RuleVersion + RequestId + sorted context data`. Sorting the dictionary ensures key order does not affect the hash. The same inputs always produce the same hash. A hash mismatch during replay signals input drift or a code change.

---

### Q: What is `TimeProvider` and why is it injected?

**Short answer:** `TimeProvider` is a .NET 8+ abstraction over the system clock. Injecting it rather than calling `DateTime.Now` directly makes the evaluator testable and deterministic—tests can set a fixed time.

**Elaboration:** This is the same pattern as the Go `Clock` interface. Without it, any call to `DateTime.UtcNow` inside evaluation would produce a different hash on replay, breaking the determinism guarantee. By injecting `TimeProvider`, the evaluation is pure: the same rule + context + time = identical output every time.

---

### Q: What is the `Rule` model? Why is it a `sealed record`?

**Short answer:** A `sealed record` in C# is an immutable value type with structural equality and no inheritance. It perfectly models a rule version: once created, it cannot change.

**Elaboration:**
```csharp
public sealed record Rule {
    public required string Id { get; init; }
    public required int Version { get; init; }
    public required IReadOnlyList<RuleCondition> Conditions { get; init; }
    ...
}
```
`init`-only properties mean fields can only be set during construction. `IReadOnlyList` prevents mutation of conditions. This enforces the immutability invariant at the type level—a rule version is a fact, not a mutable object.

---

### Q: How are conditions evaluated? What operators are supported?

**Short answer:** Conditions are evaluated recursively. Each condition specifies a field, an operator, and a value. Logical operators (`And`, `Or`, `Not`) combine results.

**Elaboration:** The `EvaluateConditions` method iterates conditions and combines results left-to-right using the preceding condition's `LogicalOperator`. Supported comparison operators: `Equals`, `NotEquals`, `GreaterThan`, `GreaterThanOrEqual`, `LessThan`, `LessThanOrEqual`, `Contains`, `NotContains`, `StartsWith`, `EndsWith`, `In`, `NotIn`, `Regex`. Comparisons are type-aware: numeric fields use decimal comparison, date fields use `DateTimeOffset` comparison, everything else uses string comparison.

---

### Q: How does rule versioning work?

**Short answer:** Every update creates a new MongoDB document with an incremented version number. Old versions are never deleted.

**Elaboration:** The `rules` collection has a unique index on `{RuleId: 1, Version: -1}`. When a rule is updated, the service inserts a new document with `Version + 1` rather than updating the existing one. This means historical decisions are always traceable to their exact rule version, and replay can re-evaluate against the original version, a newer version, or any version in between.

---

### Q: What is the architecture of the Decision Engine project?

**Short answer:** Clean architecture with four layers: `Domain` (models, interfaces), `Application` (services, evaluator, DTOs, validators), `Infrastructure` (MongoDB repos, HTTP client), `Api` (endpoints, middleware).

**Elaboration:**
- `Domain` has zero external dependencies—pure C# records, interfaces, and exceptions.
- `Application` depends only on `Domain`—business logic, evaluation, and mapping.
- `Infrastructure` depends on `Application` and `Domain`—MongoDB drivers, HTTP clients.
- `Api` wires everything together via DI in `Program.cs` using minimal APIs.
This layering means the evaluator and service logic can be unit-tested without any database or network dependency.

---

### Q: What is an `EvaluationContext`?

**Short answer:** The runtime input to the evaluator: a dictionary of key-value pairs extracted from the event payload plus metadata, and a `RequestId` for hash correlation.

**Elaboration:** The context is the bridge between an event (opaque JSON) and the rule evaluator (typed conditions). The Decision Engine extracts relevant fields from the event payload into `context.Data` (a `Dictionary<string, object>`). The evaluator then looks up each `condition.Field` in this dictionary and compares it against `condition.Value`.

---

## 5. Query & Audit Service (Go)

---

### Q: What does the Query & Audit service do?

**Short answer:** It is a read-optimized projection layer. It indexes events and decisions from the other services into Elasticsearch and provides temporal search, decision diff, audit trail, and correlation APIs.

---

### Q: What are the sync workers and how do they work?

**Short answer:** Background goroutines that periodically poll Event Timeline and Decision Engine, page through new records using a cursor, and bulk-index them into Elasticsearch. They use exponential backoff on failure.

**Elaboration (code):**
```go
// sync/worker.go
func (w *Worker) run(ctx context.Context, cursor string) {
    backoff := w.interval
    for {
        newCursor, err := w.syncFunc(ctx, cursor)
        if err != nil {
            backoff = min(backoff*2, w.maxBackoff) // exponential backoff
            continue
        }
        cursor = newCursor
        backoff = w.interval // reset on success
    }
}
```
The cursor (timestamp or sequence) is persisted in Elasticsearch (`aevum-sync-state` index) so that on restart the worker resumes from where it left off, not from the beginning.

---

### Q: What is the `DiffEngine` and when would you use it?

**Short answer:** It compares two sets of decisions at two points in time (T1 and T2) and reports which decisions were added, removed, or had their status/rule version changed.

**Elaboration:** Use case: you want to promote a new rule version to production. Before doing so, run a "what-if" replay of the last 30 days of events against the new rule version. The `DiffEngine` compares the resulting simulated decisions against the original historical decisions. If too many decisions change from `Approved` to `Rejected`, you can block the promotion.

---

### Q: How does temporal search work?

**Short answer:** Events and decisions are indexed into Elasticsearch with timestamps. The temporal query API accepts `from` and `to` timestamps, translates them into Elasticsearch `range` filters, and returns matching documents.

**Elaboration (code):**
```go
// search/temporal_query.go
filters = append(filters, map[string]interface{}{
    "range": map[string]interface{}{
        "occurred_at": map[string]interface{}{ "gte": from, "lte": to },
    },
})
```
Elasticsearch is well-suited for this pattern: range queries on date fields use inverted indices and are fast even over millions of documents. The service also supports filtering by `stream_id` and `event_type`.

---

### Q: What is an audit trail in Aevum?

**Short answer:** For a given decision, the audit trail links: the original event (with its stream, payload, and metadata) → the rule version used → each condition evaluated and its result → the final status and output hash.

**Elaboration:** The `AuditBuilder` in `search/audit_builder.go` assembles this chain by fetching the indexed decision from Elasticsearch, then correlating it with the event (by `event_id`) and the rule (by `rule_id` + `rule_version`). The result is a structured trace that can answer: *"Why was this payment rejected on 12 February?"*

---

### Q: Why use Elasticsearch for search instead of DynamoDB or MongoDB?

**Short answer:** DynamoDB does not support full-text search or complex temporal range queries across multiple fields efficiently. MongoDB can do it but is not optimized for analytics-scale aggregations. Elasticsearch is purpose-built for this.

**Elaboration:** The search requirements—multi-field text search, time-range filtering, correlation across `stream_id`, `rule_id`, `event_type`—are exactly what Elasticsearch inverted indices and BKD trees for numeric/date ranges are designed for. The trade-off is an eventual-consistency sync pipeline, which is acceptable for audit/search workloads where a few seconds of lag is fine.

---

## 6. Frontend (Vue 3 + TypeScript)

---

### Q: What does the frontend do and why Vue?

**Short answer:** The frontend provides the operator and engineer UI: timeline viewer, decision inspector, replay console, and audit trail visualization. Vue was chosen for its Composition API, strong TypeScript support, and Pinia for reactive state management.

**Elaboration:** Vue 3's Composition API gives full TypeScript inference in components, Pinia is lighter than Vuex and TypeScript-native, and Vite provides fast HMR for developer productivity. The choice reflects a pragmatic preference for a modern, well-typed frontend stack without the complexity overhead of a full React ecosystem.

---

### Q: What are the key pages/views?

**Short answer:** Dashboard (metrics overview), Events/Timeline (event stream per stream), Decisions (list + detail with trace), Rules (create/update/view rules), Replay (submit + watch live replay), Audit (audit trail + causal chain graph).

---

### Q: How does the frontend communicate with the backend?

**Short answer:** Through three typed API client modules (`eventTimeline.ts`, `decisionEngine.ts`, `queryAudit.ts`) backed by a shared Axios-based `client.ts`. A `normalizers.ts` layer translates backend snake_case to frontend camelCase.

**Elaboration:** All API calls go through a central `client.ts` that applies base URLs, auth headers, and error interceptors. Each service has its own module exporting typed functions (e.g., `getDecisions()`, `createRule()`). The `types.ts` file defines shared TypeScript interfaces mirroring the OpenAPI schemas, giving end-to-end type safety from backend response to UI component props.

---

### Q: How does replay streaming work in the UI?

**Short answer:** The frontend opens a Server-Sent Events (SSE) connection to the Event Timeline admin API. As the replay engine emits events through a Go channel, they are flushed to the SSE stream and rendered progressively in `ReplayEventFeed.vue`.

---

## 7. DevOps, Infrastructure & CI/CD

---

### Q: How is the local development environment set up?

**Short answer:** `docker compose up -d --build` starts the full stack: all three services, MongoDB, DynamoDB Local, Elasticsearch, and the Vue frontend. A `seed-data` container populates test data automatically.

---

### Q: What does the Pulumi code provision?

**Short answer:** A full AWS production environment: VPC + subnets, EKS cluster, DynamoDB table, DocumentDB cluster, OpenSearch domain, ECR registries, Lambda (event fanout), SQS queues, CloudFront + S3 (frontend CDN), Route53 DNS, IAM roles with IRSA, Prometheus and Grafana on EKS.

**Elaboration:** The `index.ts` entry point orchestrates the entire infrastructure. Component functions (`createVpc`, `createEksCluster`, etc.) are in `components/`. Configuration per environment (dev, staging, prod) lives in `Pulumi.dev.yaml`, `Pulumi.staging.yaml`, `Pulumi.prod.yaml`. IRSA (IAM Roles for Service Accounts) gives each pod least-privilege AWS access without storing credentials in code or Secrets.

---

### Q: What is IRSA and why is it used?

**Short answer:** IRSA = IAM Roles for Service Accounts. It binds a Kubernetes service account to an AWS IAM role via OIDC federation. Each pod gets temporary credentials scoped to exactly what it needs—no shared credentials, no IAM user keys in environment variables.

---

### Q: What is ArgoCD's ApplicationSet?

**Short answer:** An ApplicationSet generates one ArgoCD Application per service directory from a single template. Adding a new service to the `devops/argocd/apps/` directory automatically creates a deployment pipeline without touching the ApplicationSet.

---

### Q: What does the CI pipeline validate?

**Short answer:** Docker image builds for all services (GitHub Actions `docker-build.yml`), CodeQL static analysis for security (`codeql.yml`), and GitLab CI runs service-specific lint/test/build steps.

---

### Q: Why both GitHub Actions and GitLab CI?

**Short answer:** GitHub Actions is the primary CI for the public monorepo (Docker build validation, CodeQL security scanning). The `.gitlab-ci.yml` serves as an alternative/legacy pipeline definition demonstrating multi-platform portability.

---

## 8. Observability

---

### Q: How is observability implemented?

**Short answer:** All three services emit distributed traces and metrics via OpenTelemetry (OTLP protocol) to an OTel Collector, which scrapes into Prometheus. Grafana visualizes Prometheus metrics. Logs are structured JSON with correlation IDs.

**Elaboration (ADR-007):** The standard is vendor-neutral. Traces carry `trace_id` and `span_id` that correlate across all three services—you can see a single request as it flows through Event Timeline → Decision Engine → Query & Audit in one Jaeger/Tempo view. Prometheus metrics track: ingestion rate, replay duration, active replays, decision evaluation latency, error rates.

---

### Q: What metrics does the Event Timeline service expose?

**Short answer:** `ingest_events_total` (counter by stream/type/status), `ingestion_duration_seconds` (histogram), `replay_events_total` (counter), `active_replays` (gauge), `replay_duration_seconds` (histogram).

---

### Q: How are correlation IDs propagated?

**Short answer:** A `request_id` middleware generates a UUID per request and sets it in both the response header (`X-Request-ID`) and the structured log context. Downstream service calls pass it as a header, threading a single logical operation through all service logs.

---

## 9. Security

---

### Q: How is authentication handled?

**Short answer:** JWT (HS256) on public APIs. Admin APIs are on a separate internal port and not exposed to the public network.

**Elaboration (code):**
```go
// middleware/auth.go
token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
    if t.Method.Alg() != jwt.SigningMethodHS256.Alg() { return nil, errors.New("unexpected algorithm") }
    return []byte(secret), nil
})
```
The middleware validates algorithm (prevents `alg:none` attacks), signature, and required claims (`iss`, `sub`, `exp`, `iat`). An invalid or missing token returns 401 immediately.

---

### Q: Where are secrets stored?

**Short answer:** Never in code. Locally, environment variables via Docker Compose. In Kubernetes, Kubernetes Secrets (referenced in Helm chart values). In production AWS, AWS Secrets Manager with IRSA-based access.

---

### Q: What security features protect the data stores?

**Short answer:** DynamoDB: AWS-managed encryption at rest, IAM least-privilege via IRSA. OpenSearch: VPC-confined, security group restricted. DocumentDB/MongoDB: TLS enforced, VPC private subnets only.

---

### Q: How do you prevent `alg:none` JWT attacks?

**Short answer:** The middleware explicitly checks `token.Method.Alg() != jwt.SigningMethodHS256.Alg()` and rejects any token signed with a different (or no) algorithm before even attempting signature verification.

---

## 10. Testing

---

### Q: What types of tests exist?

**Short answer:** Unit tests for all services (domain logic, evaluator, validators, mappers), integration tests for Decision Engine and Event Timeline (using Docker containers), and Go table-driven tests throughout.

**Elaboration:**
- Decision Engine: `xUnit` tests in `Aevum.DecisionEngine.Application.Tests` and `Aevum.DecisionEngine.Domain.Tests` cover `DeterministicEvaluator`, `EvaluationService`, `RuleManagementService`, validators, and mappers. Integration tests in `Aevum.DecisionEngine.Integration.Tests` spin up a real MongoDB instance via Testcontainers.
- Event Timeline: Go table-driven unit tests for ingest, replay, storage, and API layers. Integration tests use a local DynamoDB container via `testhelpers/`.
- Query & Audit: Go tests for search engine, diff engine, sync worker, temporal query, and storage layers.

---

### Q: How is the `DeterministicEvaluator` tested for determinism?

**Short answer:** Tests inject a fixed `TimeProvider`, pass identical rule and context inputs, and assert that the output hash is identical across multiple invocations and across both original and replay paths.

---

### Q: What testing gaps exist?

**Short answer:** No Playwright end-to-end tests for the frontend, no chaos/resilience tests, no contract tests between services.

**Elaboration:** These are explicitly acknowledged in the README's "What I Would Do Differently" section.

---

## 11. Data Model Deep-Dives

---

### Q: Explain the DynamoDB single-table design.

**Short answer:** One table, multiple entity types distinguished by PK/SK patterns. Events, sequence guards, and idempotency locks all live in the same table, using key prefixes to differentiate them.

**Elaboration:**
| PK | SK | Purpose |
|---|---|---|
| `evt_01J3ZQ8A9R` | `EVENT#account-123#00000000000000000042` | Event record |
| `SEQ#account-123` | `42` | Sequence guard (prevents duplicate sequence) |
| `IDEMP#account-123#idem-87ab` | `LOCK` | Idempotency lock |
GSI1 (`GSI1PK=StreamID`, `GSI1SK=SequenceNumber`) supports ordered stream queries. GSI2 (`GSI2PK=idempotency composite key`) supports idempotency lookups.

---

### Q: Why is the sequence guard a separate item in the same transaction?

**Short answer:** It uses DynamoDB's conditional write guarantee to atomically claim a sequence number. If two concurrent writes attempt the same sequence, the transaction for one succeeds and the other fails with `ConditionalCheckFailed`.

---

### Q: How does MongoDB rule versioning work concretely?

**Short answer:** New rule version = new document with `Version + 1`. The unique index `{RuleId: 1, Version: -1}` prevents duplicate versions. Old versions remain queryable forever.

---

### Q: What is `aevum-sync-state` in Elasticsearch?

**Short answer:** A cursor persistence index. Sync workers store their last successfully processed timestamp or sequence. On restart, they resume from this cursor rather than re-indexing everything.

---

## 12. What Would You Do Differently?

These are explicitly stated in the README. Know them cold—interviewers often ask.

---

### Q: What would you add or change if you had more time?

**Short answer (from README):**

1. **gRPC for internal service contracts** — stronger interface evolution and binary efficiency compared to HTTP/JSON.
2. **Schema registry** — enforce event payload compatibility; prevent breaking changes to event structure from reaching downstream services undetected.
3. **Playwright E2E tests** — cover critical frontend and replay workflows with browser automation.
4. **Formal CQRS boundaries** — dedicated read models for audit/search, separating them cleanly from write paths.
5. **Chaos experiments** — pod kill, packet loss, and dependency latency tests to validate resilience claims empirically rather than theoretically.

---

### Q: Why didn't you use Kafka?

**Short answer:** Kafka was considered and rejected (ADR-004). DynamoDB Streams + Lambda + SQS achieved the same fanout semantics with less operational complexity for the current platform scope.

**Elaboration:** Kafka adds a broker to operate, version, monitor, and scale. For three services with bounded throughput requirements, DynamoDB Streams provides partition-ordered delivery and Lambda provides serverless processing. If the platform scaled to thousands of streams per second with complex consumer group requirements, Kafka would become worth the operational investment.

---

### Q: Why didn't you use a full event sourcing framework like Axon?

**Short answer:** Framework lock-in and cognitive overhead without proportional benefit at this scale (ADR-005). A custom implementation is smaller, more explicit, and more demonstrable.

**Elaboration:** Full frameworks like Axon or EventStoreDB are powerful but bring opinionated patterns, version constraints, and abstraction layers that can obscure understanding. Building a focused custom replay engine demonstrates the core concepts clearly and keeps the platform portable.

---

## 13. Additional Deep-Dive Questions

These are questions that probe corners of the codebase not covered in the main sections above. A thorough interviewer who has read the code will ask these.

---

### Q: What format are event IDs? Why not UUID?

**Short answer:** Event IDs are ULIDs (Universally Unique Lexicographically Sortable Identifiers), generated with the `oklog/ulid` library and seeded from the injected `Clock`.

**Elaboration (code):**
```go
// pkg/identifier/id.go
func (g *ULIDGenerator) New(t time.Time) (string, error) {
    id, err := ulid.New(ulid.Timestamp(t), g.entropy)
    ...
    return id.String(), nil
}
```
ULIDs embed a millisecond timestamp in the first 10 bytes and 16 bytes of randomness in the rest. Unlike UUIDs (version 4), ULIDs are *lexicographically sortable*—newer events naturally sort after older ones in string comparisons, which aligns well with DynamoDB's lexicographic SK ordering. They are also URL-safe, making them safe for path parameters without encoding.

---

### Q: How does rate limiting work in the public API?

**Short answer:** Per-IP token bucket using `golang.org/x/time/rate`. Each IP gets its own limiter. Visitors idle for more than 3 minutes are cleaned up by a background goroutine.

**Elaboration (code):**
```go
// middleware/ratelimit.go
visitors := map[string]*visitor{}
// ...
if !limiter.Allow() {
    c.Header("Retry-After", "1")
    httputil.TooManyRequests(c, "rate_limited", ...)
    c.Abort()
}
```
Configurable via `RATE_LIMIT_RATE` (requests/sec) and `RATE_LIMIT_BURST` (default: 50 r/s, 100 burst from Helm values). The map is protected by a mutex. A separate goroutine runs every minute to evict stale entries and prevent unbounded memory growth.

---

### Q: How does cursor-based pagination work in the stream endpoint?

**Short answer:** The stream handler encodes `{streamId}:{sequence}:{direction}` as a base64 string. The client passes this opaque token back as `?cursor=...` to resume exactly where the last page ended.

**Elaboration (code):**
```go
// domain/cursor.go
func (c Cursor) Encode() string {
    raw := fmt.Sprintf("%s:%d:%s", c.StreamID, c.Sequence, c.Direction)
    return base64.StdEncoding.EncodeToString([]byte(raw))
}
```
The handler validates that the cursor's `StreamID` matches the request path—preventing cross-stream cursor misuse. If `has_more` is false, `next_cursor` is empty and the client knows pagination is exhausted. Direction can be `forward` (ascending sequence) or `backward` (descending from latest).

---

### Q: What is the batch ingest endpoint and what are its constraints?

**Short answer:** `POST /api/v1/events/batch` accepts 1–25 events in a single call. Each event is validated and ingested independently; per-item results include `status` (`created`/`duplicate`/`invalid`/`error`) and an optional `error` message.

**Elaboration:** The 25-item limit mirrors DynamoDB's `BatchWriteItem` API limit and prevents unbounded request bodies. Each event goes through the same full ingest path—idempotency check, sequence assignment, transactional write—so the batch is not atomic: some items may succeed and others fail independently. The response contains an array of `BatchResult` objects so the caller can inspect per-item outcomes.

---

### Q: How does the Decision Engine prevent evaluating the same request twice?

**Short answer:** Before evaluating, `EvaluationService` looks up `RequestId` in the decisions collection. If found, the existing decision is returned immediately—no re-evaluation.

**Elaboration (code):**
```csharp
// EvaluationService.cs
var existingDecision = await _decisionRepository.GetByRequestIdAsync(context.RequestId, ...);
if (existingDecision is not null) return existingDecision;
// ... evaluate only if not already exists
```
The `RequestId` has a unique index in MongoDB (`{ RequestId: 1 }`, unique). There is also a race-condition guard: if two concurrent requests with the same `RequestId` race past the initial check, only one `InsertOneAsync` will succeed; the other catches the exception and retries the lookup, finding the winner's result. This makes decision evaluation idempotent at the service boundary.

---

### Q: What happens after a decision is evaluated—does the Decision Engine do anything else?

**Short answer:** Yes. It fires-and-forgets an event back to the Event Timeline Service (`decision.evaluated` event type), making the decision itself observable in the event stream.

**Elaboration (code):**
```csharp
// EvaluationService.cs
_ = PublishTimelineEventAsync(savedDecision);

private async Task PublishTimelineEventAsync(Decision decision) {
    try {
        await _eventTimelineClient.IngestEventAsync(
            streamId: $"decision-{decision.Id}",
            eventType: "decision.evaluated", ...);
    } catch { /* silently fail */ }
}
```
This is a best-effort fire-and-forget: the task is not awaited. If the Event Timeline is unavailable, the failure is swallowed. The rationale is that publishing the decision event is a side effect—not critical to the evaluation response—and failing it should not fail the caller.

---

### Q: How does the Decision Engine handle HTTP failures when calling Event Timeline?

**Short answer:** Polly policies: 3 retries with exponential backoff (100ms base), then a circuit breaker that opens after 5 consecutive failures and stays open for 30 seconds.

**Elaboration (code):**
```csharp
// ServiceExtensions.cs
private static IAsyncPolicy<HttpResponseMessage> GetRetryPolicy() {
    return HttpPolicyExtensions.HandleTransientHttpError()
        .WaitAndRetryAsync(3, retryAttempt => TimeSpan.FromMilliseconds(100 * Math.Pow(2, retryAttempt)));
}
private static IAsyncPolicy<HttpResponseMessage> GetCircuitBreakerPolicy() {
    return HttpPolicyExtensions.HandleTransientHttpError()
        .CircuitBreakerAsync(5, TimeSpan.FromSeconds(30));
}
```
"Transient" means 5xx responses or network errors. The circuit breaker prevents the Decision Engine from hammering a failing Event Timeline during recovery. When `EventTimeline:BaseUrl` is not configured at all, a `NoOpEventTimelineClient` is injected instead—allowing the Decision Engine to run completely standalone without any dependency on Event Timeline.

---

### Q: What are the Rule lifecycle states? How does a rule go from creation to active?

**Short answer:** Three states: `Draft` (created, not yet active), `Active` (evaluated against incoming events), `Inactive` (deactivated, no longer evaluated).

**Elaboration:**
1. `POST /api/v1/rules` → creates at `Version: 1, Status: Draft`.
2. `POST /api/v1/rules/{id}/activate` → transitions to `Active`.
3. `POST /api/v1/rules/{id}/deactivate` → transitions to `Inactive`.
4. `PUT /api/v1/rules/{id}` → creates a new version (`Version + 1`) with the updated fields; the new version's status is whatever was passed in the request.

Active rules are loaded by `GetActiveRulesAsync`, which queries MongoDB by `Status = Active` sorted by descending `Priority`. Only `Active` rules are candidates for evaluation.

---

### Q: What are the possible `DecisionStatus` values?

**Short answer:** `Approved` (all conditions matched), `Rejected` (conditions did not match), `Error` (evaluation threw an exception).

**Elaboration:** The evaluator sets `Approved` when `isMatch = true` and `Rejected` when `isMatch = false`. `Error` is set in `ExceptionHandling` middleware when an `EvaluationException` is caught—the decision is still persisted with status `Error` and the error message, providing a complete audit trail even for failed evaluations.

---

### Q: What exactly does the Lambda fanout function do?

**Short answer:** It reads DynamoDB Stream records (INSERT and MODIFY events only), transforms them into SQS messages, and sends them in batches of 10 to the `event-notifications` SQS queue.

**Elaboration (code):**
```typescript
// lambda-code/fanout/index.ts
const messages = event.Records
    .filter(r => r.eventName === "INSERT" || r.eventName === "MODIFY")
    .map((r, i) => ({
        Id: `${r.eventID}-${i}`.slice(0, 80),
        MessageBody: JSON.stringify({ eventName, keys, newImage, oldImage, sequenceNumber, ... })
    }));

for (const batch of chunk(messages, 10)) {
    await sqs.send(new SendMessageBatchCommand({ QueueUrl, Entries: batch }));
}
```
`DELETE` records are intentionally ignored—events are immutable in Aevum so deletions should not occur; filtering them prevents phantom signals to consumers. SQS `SendMessageBatch` has a 10-message limit, so the Lambda chunks into groups of 10 and logs `{ sent, failed }` for observability.

---

### Q: What DynamoDB table features does the Pulumi config enable?

**Short answer:** On-demand billing (`PAY_PER_REQUEST`), DynamoDB Streams with `NEW_AND_OLD_IMAGES`, point-in-time recovery (PITR), and server-side encryption (SSE).

**Elaboration:**
- **On-demand billing**: no capacity planning required; scales automatically to request volume.
- **Streams with `NEW_AND_OLD_IMAGES`**: the Lambda fanout function receives both the before and after state of each item, enabling change-data-capture semantics.
- **PITR**: allows table restore to any second within the last 35 days—disaster recovery without explicit snapshots.
- **SSE**: data at rest encrypted with AWS-managed keys (no customer key management overhead).

---

### Q: Explain the Helm chart's NetworkPolicy and why it matters.

**Short answer:** The NetworkPolicy restricts ingress to the pod to traffic from the same `project: aevum` namespace only, and allows all egress. This prevents external pods or other namespaces from directly hitting the service, enforcing network micro-segmentation.

**Elaboration:**
```yaml
# templates/networkpolicy.yaml
ingress:
  - from:
      - namespaceSelector:
          matchLabels:
            project: aevum
    ports:
      - protocol: TCP
        port: 8080
```
Without this, any pod in any namespace in the cluster could reach the service directly, bypassing ingress/service mesh controls. With it, only pods in a namespace labelled `project: aevum` (i.e., other Aevum services, the ingress controller, and monitoring) can connect.

---

### Q: What does the HPA in the Helm chart do and what are its settings?

**Short answer:** The HorizontalPodAutoscaler scales the deployment between 2 and 10 replicas based on CPU utilization (target 70%) and memory utilization (target 80%).

**Elaboration:**
```yaml
autoscaling:
  minReplicas: 2    # Never below 2 for HA
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70
  targetMemoryUtilizationPercentage: 80
```
The minimum of 2 replicas ensures high availability across two availability zones even at idle load. CPU 70% gives headroom before saturation. Memory at 80% is conservative to avoid OOM kills—memory scaling is less elastic than CPU so you want to scale before you hit the limit.

---

### Q: What is the `checksum/config` annotation on the Deployment?

**Short answer:** It is a SHA-256 hash of the ConfigMap content, embedded as a pod annotation. When the ConfigMap changes, the hash changes, and Kubernetes rolls the Deployment automatically.

**Elaboration:**
```yaml
# templates/deployment.yaml
annotations:
  checksum/config: {{ include (print $.Template.BasePath "/configmap.yaml") . | sha256sum }}
```
Without this annotation, updating a ConfigMap does not trigger a pod restart by default—pods continue running with stale config. This pattern is a standard Helm idiom to force a rolling restart when configuration changes, ensuring pods always run with the latest config values.

---

### Q: What is a ServiceMonitor and why is it in the Helm chart?

**Short answer:** A `ServiceMonitor` is a Prometheus Operator CRD that tells Prometheus which service to scrape, on which port, and at what interval—without touching Prometheus config files directly.

**Elaboration:**
```yaml
serviceMonitor:
  enabled: true
  interval: 15s
  path: /admin/metrics
  port: admin
```
The admin port on the Event Timeline exposes Prometheus metrics via `/admin/metrics`. The ServiceMonitor tells the Prometheus Operator to scrape that endpoint every 15 seconds. This approach keeps scrape configuration declarative and versioned in Git alongside the service—no separate Prometheus config management.

---

### Q: What Prometheus alerting rules exist and what do they fire on?

**Short answer:** Seven rules across four groups. The most critical: `HighIngestionErrorRate` (error rate > 5% for 5 min → critical), `DecisionLatencyHigh` (P99 > 2s for 5 min → critical), `SearchSyncLag` (lag > 5 min → warning).

**Full list:**
| Alert | Condition | Severity |
|---|---|---|
| `HighIngestionErrorRate` | error rate > 5% (5 min) | critical |
| `IngestionLatencyHigh` | P95 ingestion > 1s (5 min) | warning |
| `DecisionEvaluationErrors` | any evaluation errors (10 min) | warning |
| `DecisionLatencyHigh` | P99 evaluation > 2s (5 min) | critical |
| `SearchSyncLag` | sync lag > 300s (10 min) | warning |
| `PodRestartingFrequently` | > 3 restarts in 1h | warning |
| `HighMemoryUsage` | memory > 90% of limit (5 min) | critical |

---

### Q: What Prometheus recording rules exist and why?

**Short answer:** Eight pre-computed metrics (rates and percentiles at P50/P95/P99) that make dashboards faster by avoiding expensive on-the-fly histogram calculations.

**Elaboration:**
```yaml
# recording-rules.yaml
- record: aevum:ingestion_latency_p95:5m
  expr: histogram_quantile(0.95, rate(aevum_ingestion_duration_seconds_bucket[5m]))
```
Recording rules persist pre-computed query results as new metric series on a 30-second interval. Grafana dashboards query the pre-computed `aevum:*` series instead of running expensive `histogram_quantile(rate(...))` expressions on every panel load—especially important at scale when the cardinality of the underlying series is high.

---

### Q: Walk me through the OTel Collector configuration.

**Short answer:** Receives OTLP on gRPC (4317) and HTTP (4318), applies batch and memory-limiter processors, exports metrics to a Prometheus endpoint on 8888 and traces to the `debug` exporter.

**Elaboration:**
```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc: {endpoint: 0.0.0.0:4317}
      http: {endpoint: 0.0.0.0:4318}
processors:
  batch: {timeout: 5s, send_batch_size: 1024}
  memory_limiter: {limit_mib: 512, spike_limit_mib: 128}
exporters:
  prometheus: {endpoint: 0.0.0.0:8888, namespace: aevum}
```
Services send OTLP traces and metrics to the Collector. The `batch` processor amortizes network calls by grouping signals into 1024-item batches or 5-second windows. The `memory_limiter` protects the Collector from OOM under burst load by dropping data when memory exceeds 512 MiB. Prometheus scrapes the `:8888` endpoint. Traces go to `debug` (logged) locally; in production they would be exported to Tempo or Jaeger.

---

### Q: What is a CorrelationQuery and how is it different from the DiffEngine?

**Short answer:** `CorrelationQuery` finds all events and decisions that share a common identifier (event_id, rule_id, stream_id, etc.) at a point in time. `DiffEngine` compares two *sets* of decisions at two points in time.

**Elaboration:** Use `CorrelationQuery` to answer: *"Show me everything related to event `evt_01J3` — which decisions were triggered, in which streams, under which rules?"* Use `DiffEngine` to answer: *"Between last week and today, which decisions changed status?"* The `CorrelationQuery` is a multi-field `bool.filter` query against both `aevum-events` and `aevum-decisions` indices simultaneously. It supports filtering by `event_id`, `decision_id`, `rule_id`, `rule_version`, `stream_id`, and `event_type` in any combination.

---

### Q: How does the BulkIndexer work?

**Short answer:** It buffers documents in memory and flushes them to Elasticsearch using the `_bulk` API when the buffer reaches `batchSize * 2` items (one metadata line + one document line per doc). A forced flush is called at the end of each sync cycle.

**Elaboration:**
```go
// indexer/bulk_indexer.go
func (bi *BulkIndexer) IndexDocument(ctx context.Context, index, id string, doc interface{}) error {
    bi.addToBatch(index, id, doc)
    if len(bi.buffer) >= bi.batchSize*2 {
        return bi.Flush(ctx)
    }
    return nil
}
```
The buffer holds alternating metadata/document pairs as required by the ES bulk API format. Flushing serialises each pair as newline-delimited JSON and sends one HTTP request. After a successful flush the buffer is cleared. Default `batchSize` is 500 (from config), so the buffer flushes at 1000 lines.

---

### Q: How does the Query & Audit service authenticate to the Event Timeline service?

**Short answer:** It constructs a JWT from scratch—header + payload base64-encoded, HMAC-SHA256 signed with the shared `EVENT_TIMELINE_JWT_SECRET` environment variable.

**Elaboration (code):**
```go
// clients/event_timeline_client.go
header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"iss":"query-audit","sub":"sync-worker","iat":%d,"exp":%d}`, now, exp)))
h := hmac.New(sha256.New, []byte(c.jwtSecret))
h.Write([]byte(header + "." + payload))
sig := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
req.Header.Set("Authorization", "Bearer "+header+"."+payload+"."+sig)
```
If `EVENT_TIMELINE_JWT_SECRET` is empty, no auth header is added—useful in local development. In production the secret is injected via Kubernetes Secret. This is the same `JWTAuth` middleware the Event Timeline validates: `iss=query-audit`, `sub=sync-worker`, 1-hour TTL.

---

### Q: What does the admin replay endpoint actually return—is it SSE?

**Short answer:** No. The current `TriggerReplay` handler is a simple JSON response, not SSE. It collects all replayed events synchronously and returns `{ "status": "completed", "events_replayed": count }`.

**Elaboration (code):**
```go
// handlers/admin/replay_handler.go
eventsCh, errCh := h.engine.Replay(ctx, req)
count := 0
for range eventsCh {
    count++
}
return c.JSON(200, map[string]any{"status": "completed", "events_replayed": count})
```
This means replay is synchronous and blocking until the full stream is consumed. For the frontend's real-time streaming experience, the Vue `ReplayEventFeed.vue` component uses SSE differently—the frontend subscribes to an SSE endpoint that the replay handler would need to be extended to support. The current admin replay is a functional backend integration point, not yet wired for SSE push. This is a known evolution point.

---

### Q: What happens when you delete a rule?

**Short answer:** All versions of that rule are deleted from MongoDB (`DeleteManyAsync` by `RuleId`).

**Elaboration (code):**
```csharp
// MongoDbRuleRepository.cs
public async Task DeleteAsync(string id, CancellationToken cancellationToken = default)
{
    var filter = Builders<RuleDocument>.Filter.Eq(r => r.RuleId, id);
    await _collection.DeleteManyAsync(filter, cancellationToken);
}
```
This is a soft delete–free design: deletion removes all version history. This trade-off is intentional—rules are versioned in the repository, and decisions already store the `RuleId + RuleVersion` they were evaluated against, so historical audit integrity is preserved in the `decisions` collection even if the rule document is deleted. If history preservation is required, the `DeleteAsync` could be changed to a status update to `Deleted` instead.

---

### Q: Why does .NET use Minimal APIs instead of MVC controllers?

**Short answer:** Minimal APIs are lightweight, have lower overhead, are easier to read, and align with modern .NET idioms. They avoid the ceremony of controller classes, base class inheritance, and attribute routing.

**Elaboration:**
```csharp
// Program.cs
var apiV1 = app.MapGroup("/api/v1");
apiV1.MapGroup("/rules").MapRuleEndpoints();
apiV1.MapGroup("/decisions").MapDecisionEndpoints();
```
With Minimal APIs, each endpoint is a function—no `ControllerBase` inheritance, no `[HttpPost]` attributes. `RouteGroupBuilder` provides grouping and shared middleware (OpenAPI tags, response type declarations). The trade-off is that complex features like model binding conventions are less automatic, but for a platform with explicit DTOs and FluentValidation, this is fine.

---

### Q: What do the Pulumi infrastructure tests cover?

**Short answer:** Unit tests for the five most critical IaC components: DynamoDB table configuration, EKS cluster settings, IAM roles and policies, VPC/subnet layout, and resource naming conventions.

**Elaboration:** The tests in `devops/pulumi/tests/` use Pulumi's testing SDK (`@pulumi/pulumi/testing`) to mock the cloud providers and assert that the resources are configured correctly—for example, that the DynamoDB table has streams enabled, that the EKS cluster is in the right VPC, and that IAM policies follow least-privilege. This catches infrastructure misconfiguration before `pulumi up` is run, not after.

---

### Q: What are the five Grafana dashboards?

**Short answer:** `aevum-overview` (platform-wide KPIs), `event-timeline` (ingestion rates, latencies, replay metrics), `decision-engine` (evaluation rates, latencies, rule hit rates), `query-audit` (sync lag, search latencies), `infrastructure` (pod restarts, memory, CPU across all services).

**Elaboration:** Each dashboard is a JSON definition in `devops/monitoring/grafana/dashboards/`. They are provisioned automatically via the `dashboard-providers.yaml` config, which tells Grafana to load all JSON files from a directory. This means dashboards are version-controlled and deployed as code—no manual clicking required. The `prometheus.yaml` datasource file wires Grafana to the Prometheus instance on startup.

---

### Q: What is Kustomize and how is it used here?

**Short answer:** Kustomize is a Kubernetes config management tool that patches base YAML manifests with environment-specific overlays. The `devops/k8s/overlays/sit` overlay patches the base manifests for the SIT (System Integration Testing) environment.

**Elaboration:**
```
devops/k8s/
├── base/          # Base manifests: services, deployments, configmaps
└── overlays/
    └── sit/       # SIT-specific patches (replicas, resource limits, etc.)
        └── kustomization.yaml
```
ArgoCD's `aevum-sit.yaml` application points to this overlay, so ArgoCD applies the SIT-patched manifests to the SIT cluster. This is different from the Helm-based production deployment—the base K8s manifests are simpler, developer-friendly, and used for integration testing rather than production operations.

---

### Q: What structured logging library is used in the Go services?

**Short answer:** `log/slog` from the Go standard library (added in Go 1.21). Both the Event Timeline and Query & Audit services use it directly.

**Elaboration:**
```go
slog.Error("ingest failed", slog.String("error", err.Error()))
slog.Info("sync successful", slog.String("service", w.serviceName))
```
`slog` outputs structured key-value JSON by default, which is compatible with log aggregation systems (e.g., CloudWatch Logs, Loki). The Decision Engine uses Serilog (a .NET structured logging library) configured via `builder.AddSerilog()`. Both produce JSON-formatted logs with correlation IDs, making cross-service log correlation possible.

---

### Q: How does the Gin recovery middleware work?

**Short answer:** It wraps all handlers with a deferred recover, logs the panic with the request ID, and returns a 500 Internal Server Error instead of crashing the process.

**Elaboration (code):**
```go
// middleware/recovery.go
func Recovery(logger *slog.Logger) gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered any) {
        logger.Error("panic recovered", slog.Any("error", recovered),
            slog.String("request_id", c.GetString(RequestIDContextKey)))
        httputil.Internal(c, "internal_error", "internal server error")
    })
}
```
Without recovery middleware, a panic in any handler goroutine would crash the entire server process. With it, panics are isolated per request. The `RequestIDContextKey` is included in the log so the panic can be correlated with the specific request that caused it.

---

### Q: What does `replay/options.go` do and what is the default page size?

**Short answer:** `NewOptions` builds the replay filter configuration—time window, event type set, and page size. Default page size is 100 if the caller passes 0 or a negative value.

**Elaboration:** The `event_types` field in the replay request is a slice of strings; `NewOptions` converts it to a `map[string]struct{}` for O(1) lookups in the replay engine's `matchesType` filter. The page size controls how many events are fetched from DynamoDB per `QueryByStream` call during replay—balancing DynamoDB read unit consumption against memory usage.

---

### Q: What fields does the `admin/replay` endpoint actually accept?

**Short answer:** Five fields: `stream_id` (required), `from` (ISO8601 timestamp), `to` (ISO8601 timestamp), `event_types` (optional string array to filter), and `page_size` (optional, default 100).

**Elaboration:** `stream_id` is the partition to replay. `from`/`to` bound the replay window. `event_types` allows replaying only specific event categories—for example, replaying only `payment_received` events to re-evaluate payment rules. `page_size` controls how many events are loaded from DynamoDB per page; lower values reduce memory footprint at the cost of more round trips.

---

### Q: What indexes does the MongoDB decisions collection create?

**Short answer:** Three indexes: unique on `RequestId` (idempotency), non-unique on `DeterministicHash` (hash lookup for replay verification), and compound on `{RuleId, RuleVersion}` (rule-scoped decision queries).

**Elaboration:**
```csharp
// MongoDbDecisionRepository.cs
var requestIdIndex = ... { Unique = true }; // prevents duplicate decision for same requestId
var hashIndex = ... DeterministicHash;      // enables hash-based deduplication queries
var ruleIndex = ... RuleId + RuleVersion;   // enables GetDecisionsByRuleId queries
```
The unique `RequestId` index is the primary idempotency guard at the database level. The `DeterministicHash` index lets the replay verification logic quickly query decisions by hash for comparison without a full collection scan.

---

## 14. Behavioral & Process Questions

European interviews often include behavioral rounds alongside technical grilling. These are specific to this project.

---

### Q: This is a solo project. How did you manage scope and avoid feature creep?

**Short answer:** I focused on the core invariants—immutable events, versioned rules, deterministic replay—and used ADRs to document and freeze decisions. Everything else (gRPC, schema registry, chaos tests) was explicitly deferred and documented in "What I Would Do Differently."

**Elaboration:** The ADR process forced me to make one decision at a time and document why alternatives were rejected. This created a defensible audit trail for every significant choice. Deferred items are honest—I didn't pretend they don't exist; they are in the README. Scope discipline came from asking: *"Does this serve the core determinism contract?"*—if not, it was deferred.

---

### Q: If someone joined your team, what would they need to know first?

**Short answer:** Read `docs/replay-model.md` first—it defines the platform's core contract. Then `docs/architecture.md` for the C4 view. Then `docker compose up` and see it running before touching code.

**Elaboration:** The replay model is the "why"—it explains the immutability, versioning, and hash verification invariants. Everything else—the code, the infrastructure, the UI—exists to serve that contract. New engineers who start with the replay model understand why rules are never mutated in place, why timestamps are injected rather than read from ambient clocks, and why the sync pipeline is eventually consistent rather than transactional.

---

### Q: How would you handle a production incident where replay results don't match original decisions?

**Short answer:** Check the hash mismatch first—it pinpoints whether the issue is input drift (different rule version/context) or code-path drift (evaluator logic changed). Then check deployment history.

**Elaboration:**
1. Query decisions by `DeterministicHash` for the affected `RequestId`.
2. If hashes differ: compare `RuleVersion` and `InputContext` between original and replay decision records.
3. If inputs match but hashes differ: a code change in `DeterministicEvaluator` altered the evaluation logic—check deployment timestamps.
4. If inputs differ: the replay was run with the wrong rule version or modified context—check the replay request parameters.
5. Mitigate: re-run replay pinning the exact original `RuleVersion` to reproduce the original decision.

---

### Q: How would you scale this platform to 10x the current volume?

**Short answer:** Event Timeline scales horizontally (stateless + DynamoDB on-demand). Decision Engine scales horizontally (stateless + MongoDB read replicas). Query & Audit scales workers independently. The bottleneck at 10x would be Elasticsearch write throughput.

**Elaboration:**
- **Event Timeline**: add replicas—HPA handles it. DynamoDB on-demand mode auto-scales.
- **Decision Engine**: add replicas—MongoDB read scaling handles rule lookups. Consider caching active rules in memory with a short TTL.
- **Query & Audit**: increase sync worker parallelism or partition sync by stream. Elasticsearch write throughput is the likely ceiling—shard scaling or tiered storage (ILM) would be needed.
- **Lambda fanout**: already serverless—scales automatically with DynamoDB stream throughput.

---

### Q: What monitoring would you look at first if ingestion latency spiked?

**Short answer:** The `HighIngestionErrorRate` and `IngestionLatencyHigh` alerts would fire first. I'd look at the P95/P99 histograms in the `event-timeline` Grafana dashboard, then DynamoDB CloudWatch metrics for throttling.

**Elaboration:** If the P95 Prometheus metric `aevum:ingestion_latency_p95:5m` is elevated, I check: (1) DynamoDB `ConsumedWriteCapacityUnits` and `ThrottledRequests`—if throttling, it's a capacity issue. (2) `ErrSequenceConflict` rate—high retry rates mean hot streams causing contention. (3) OTel traces for the specific span that is slow—is it the DynamoDB write, the idempotency lookup, or the sequence read? The recording rule pre-computes P50/P95/P99 so dashboard queries are instant.

---

### Q: Why did you choose this stack for a European engineering audience specifically?

**Short answer:** The stack is deliberate and explainable: Go and C# are both standard in European enterprise engineering teams, AWS is the dominant cloud in Europe, and Pulumi/Helm/ArgoCD represent the modern cloud-native GitOps pattern.

**Elaboration:** Go is widely used in European engineering for backend services at scale (e.g., at Wise, Zalando, N26). C# / .NET is strong in the European enterprise and financial services sector. The observability stack (OTel + Prometheus + Grafana) is the CNCF standard, language-neutral, and widely known. ArgoCD + Helm is the most common production GitOps pattern in European Kubernetes shops. The choices are not exotic—they are deliberate, defensible, and recognizable to any senior European engineer.

---

## 15. Worst-Case Rapid-Fire Questions

Questions an aggressive interviewer might fire without warning.

---

**Q: What does `matchesTimeRange` do?**
A: It filters a replayed event out if its `OccurredAt` falls before `from` or after `to`. Zero values mean "no bound."

**Q: What does `SortDictionary` in the evaluator do?**
A: It sorts the context dictionary by key before hashing it so that key insertion order doesn't affect the SHA-256 hash output.

**Q: What is `ErrSequenceConflict` vs `ErrIdempotencyConflict`?**
A: `ErrSequenceConflict` = two writers raced to claim the same sequence number; retry. `ErrIdempotencyConflict` = this exact event was already written; return the existing event.

**Q: What does `GSI2PK: "{streamId}#{idempotencyKey}"` give you?**
A: A composite key that scopes idempotency lookups to a specific stream, preventing false matches across streams with the same client-supplied key.

**Q: Why is `sealed` used on C# domain models?**
A: Prevents subclassing, which could introduce mutable derived types and break the immutability invariant.

**Q: What does the `ValidationFilter` in the Decision Engine do?**
A: It intercepts request handling before the endpoint executes, runs FluentValidation validators, and returns a 400 with field-level errors if validation fails—keeping validation out of endpoint bodies.

**Q: What is the `ExceptionHandling` middleware in the Decision Engine?**
A: A global exception handler that maps domain exceptions (`InvalidRuleException`, `RuleNotFoundException`, `EvaluationException`) to appropriate HTTP status codes (400, 404, 422) and structured error responses, preventing stack traces from leaking to callers.

**Q: What does the Helm umbrella chart do?**
A: It is a meta-chart that composes all individual service Helm charts under one release, allowing the entire platform to be deployed, upgraded, or rolled back as a single Helm operation.

**Q: What is IRSA?**
A: IAM Roles for Service Accounts. Kubernetes pods assume AWS IAM roles via OIDC federation without any AWS credentials stored in the cluster.

**Q: How many retries does `PutEvent` allow for sequence conflicts?**
A: Three retries (`retries < 3` loop in `ingest/service.go`).

**Q: What does `batchWriteMaxItems = 25` mean?**
A: DynamoDB's `BatchWriteItem` API limit is 25 items per call. The code chunks batches into groups of 25 and handles unprocessed items with up to 5 retries and progressive backoff.

**Q: What algorithm does the JWT middleware enforce?**
A: HS256. Any token with a different `alg` (including `none`) is rejected immediately to prevent algorithm substitution attacks.

**Q: What does `admin/replay` accept?**
A: `{ stream_id, from, to, event_types[], page_size }` — the stream, time window, optional type filters, and page size (default 100).

**Q: Where is the sync cursor persisted?**
A: In the `aevum-sync-state` Elasticsearch index, with the field `last_processed_at` and `cursor`.

**Q: What does the Query & Audit service's `DiffEngine.diff()` return?**
A: A `DiffResult` with three arrays: `Added` (new decision IDs in T2), `Removed` (decision IDs only in T1), `Changed` (decisions whose `status` or `rule_version` changed between T1 and T2), plus a `Summary` string.

**Q: Why is `TimeProvider` a .NET 8+ class rather than a custom interface?**
A: It is the official Microsoft abstraction for testable time in .NET. Using the standard library type means it integrates with ASP.NET and the DI container without any custom wiring.

**Q: What inputs go into the deterministic hash?**
A: `SHA-256(RuleId + RuleVersion + RequestId + sorted context data)`, all serialised as camelCase JSON. The context dictionary is key-sorted so insertion order doesn't affect the hash. A mismatch on replay indicates input drift (different rule version, different context, or a code change in the evaluator).

**Q: Why does `docker-compose.override.yml` exist?**
A: It provides local developer overrides (e.g., debug ports, hot-reload volumes, relaxed resource limits) without modifying the base `docker-compose.yml` that is also used for CI.

**Q: What is the `seed-data` container?**
A: A one-shot container that runs `devops/scripts/seed-data.sh` to populate MongoDB with active rules and decisions and DynamoDB with stream events before the UI starts. It is declared as a dependency in the Compose file so data is always present before the UI loads.

**Q: What format are event IDs and why not UUID?**
A: ULID — time-sortable, lexicographically ordered, URL-safe. Unlike UUID v4 (pure random), ULIDs embed a millisecond timestamp so they naturally sort in creation order.

**Q: What is the rate limit on the public API and how is it enforced?**
A: Per-IP token bucket via `golang.org/x/time/rate`. Default 50 r/s, 100 burst. Stale visitor entries are evicted after 3 minutes of inactivity. Returns 429 with `Retry-After: 1`.

**Q: How does cursor pagination work in the stream endpoint?**
A: `Cursor{StreamID, Sequence, Direction}` is base64-encoded as `streamId:seq:direction`. Passed as `?cursor=` query param. Handler validates the cursor's stream matches the path param to prevent cross-stream misuse.

**Q: What is the batch ingest size limit and why?**
A: 1–25 items. Matches DynamoDB's `BatchWriteItem` maximum and keeps request bodies bounded.

**Q: How does `EvaluationService` prevent evaluating the same request twice?**
A: It checks `GetByRequestIdAsync(context.RequestId)` before evaluating. If found, returns the existing decision. There is also a race-condition retry guard if two concurrent requests slip past the pre-check.

**Q: What does `PublishTimelineEventAsync` do and what happens if it fails?**
A: It fire-and-forgets a `decision.evaluated` event to the Event Timeline. Failures are silently swallowed—it is a non-critical side effect.

**Q: What Polly policies wrap the Decision Engine's HTTP client to Event Timeline?**
A: Retry 3 times with exponential backoff (100ms × 2^n), then a circuit breaker that opens after 5 failures and stays open for 30 seconds.

**Q: What are the Rule statuses and what order do they follow?**
A: Draft (created) → Active (via `/activate`) → Inactive (via `/deactivate`). Only Active rules are evaluated.

**Q: What are the three `DecisionStatus` values?**
A: `Approved` (conditions matched), `Rejected` (conditions did not match), `Error` (evaluation threw an exception).

**Q: What does the Lambda fanout filter out and why?**
A: `DELETE` records. Events are immutable in Aevum so deletions shouldn't occur; filtering prevents phantom signals.

**Q: What is the SQS batch size limit in the Lambda fanout?**
A: 10 messages per `SendMessageBatchCommand` call. The function chunks the message list accordingly.

**Q: What DynamoDB features are enabled by the Pulumi config?**
A: `PAY_PER_REQUEST` billing, `NEW_AND_OLD_IMAGES` streams, point-in-time recovery (PITR), server-side encryption (SSE).

**Q: What does the Helm NetworkPolicy restrict?**
A: Ingress is allowed only from namespaces labelled `project: aevum`. All egress is allowed.

**Q: What are the HPA replica bounds for Event Timeline?**
A: min 2, max 10. Scales on CPU 70% and memory 80%.

**Q: What triggers a pod restart when a ConfigMap changes in Helm?**
A: The `checksum/config` annotation on the Deployment—a SHA-256 of the ConfigMap content. When the ConfigMap changes, the hash changes, Kubernetes detects the annotation diff and rolls the Deployment.

**Q: What is a ServiceMonitor?**
A: A Prometheus Operator CRD that declaratively configures Prometheus scrape targets. Aevum's charts create one per service pointing at the admin port's `/admin/metrics` path every 15 seconds.

**Q: Name the seven Prometheus alerting rules.**
A: `HighIngestionErrorRate`, `IngestionLatencyHigh`, `DecisionEvaluationErrors`, `DecisionLatencyHigh`, `SearchSyncLag`, `PodRestartingFrequently`, `HighMemoryUsage`.

**Q: What ports does the OTel Collector listen on and what does it export?**
A: OTLP gRPC on 4317, OTLP HTTP on 4318. Exports metrics to Prometheus on 8888; traces to `debug` exporter (logs) locally.

**Q: What is a `CorrelationQuery` in Query & Audit?**
A: A multi-field filter query across both `aevum-events` and `aevum-decisions` indices by any combination of event_id, decision_id, rule_id, rule_version, stream_id, or event_type.

**Q: How does the BulkIndexer flush?**
A: It accumulates metadata+document pairs in a buffer and flushes via the ES `_bulk` API when the buffer reaches `batchSize * 2` items or when `Flush()` is called explicitly at the end of a sync cycle.

**Q: How does the Query & Audit service authenticate to Event Timeline?**
A: It builds a JWT inline using HMAC-SHA256 (`iss: query-audit`, `sub: sync-worker`, 1h TTL) signed with the shared `EVENT_TIMELINE_JWT_SECRET` env var.

**Q: Does the admin replay endpoint use SSE?**
A: No. The current implementation collects all replayed events and returns `{ status: "completed", events_replayed: N }` as a single JSON response. SSE streaming is a future evolution point.

**Q: What happens when you call DELETE on a rule?**
A: All versions of that rule ID are deleted from MongoDB via `DeleteManyAsync`. Historical decisions still reference the rule by ID and version for audit purposes.

**Q: What is `NoOpEventTimelineClient`?**
A: A stub `IEventTimelineClient` that returns `true` without doing anything. Injected when `EventTimeline:BaseUrl` is not configured, allowing the Decision Engine to run standalone.

**Q: How does GitLab CI detect which service to lint/test?**
A: Path-based `rules: changes:` clauses—e.g., `services/event-timeline/**`. Only changed paths trigger the corresponding job, keeping CI fast in the monorepo.

**Q: What MongoDB indexes does the decisions collection have?**
A: Unique on `RequestId` (idempotency guard), non-unique on `DeterministicHash` (hash lookup), compound on `{RuleId, RuleVersion}` (rule-scoped queries).

**Q: Why use Minimal APIs in .NET instead of MVC controllers?**
A: Less ceremony, lower overhead, cleaner endpoint-as-function model. No `ControllerBase` inheritance or attribute routing—each endpoint is an explicit function call in `RouteGroupBuilder`.

**Q: What do the Pulumi tests verify?**
A: That the IaC code produces correctly configured resources: DynamoDB streams enabled, EKS in the right VPC, IAM policies least-privilege, naming conventions followed—all without running `pulumi up`.

**Q: Name the five Grafana dashboards.**
A: `aevum-overview` (platform KPIs), `event-timeline` (ingestion/replay), `decision-engine` (evaluation), `query-audit` (sync/search), `infrastructure` (pod/memory/CPU).

**Q: What is the `log/slog` package?**
A: Go's standard library structured logger (added in Go 1.21). Both Go services use it. Outputs JSON key-value pairs, compatible with log aggregation systems.

**Q: What does the Gin Recovery middleware do?**
A: Catches panics in handlers, logs them with the request ID, and returns 500—preventing a single handler panic from crashing the server process.

**Q: What is the default replay page size?**
A: 100 events per `QueryByStream` call (set in `replay/options.go` when the caller passes 0 or negative).

**Q: What is `replay.Collect()`?**
A: A helper that drains a `<-chan domain.Event` channel into a `[]domain.Event` slice. Used in tests and non-streaming consumers of the replay engine.

**Q: What is the `dynamodb-init-job.yaml`?**
A: A Kubernetes Job that initialises the DynamoDB Local table (creates the `aevum-events` table schema) on startup. Used in the SIT environment where DynamoDB Local runs in-cluster.

**Q: What is the SIT environment?**
A: System Integration Testing—a shared environment between dev and staging. ArgoCD's `aevum-sit.yaml` deploys the platform with Kustomize overlays from `devops/k8s/overlays/sit/` for testing integration across all services.

**Q: What is `EvaluationDurationMs` on the Decision model?**
A: The wall-clock time in milliseconds that the rule evaluation took (measured by `Stopwatch`). Persisted to MongoDB and exposed in the API response for performance monitoring.

---

*End of guide. Memorize the short answers. Use the elaborations when they probe deeper. Good luck.*
