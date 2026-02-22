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
13. [Worst-Case Rapid-Fire Questions](#13-worst-case-rapid-fire-questions)

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

## 13. Worst-Case Rapid-Fire Questions

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
A: `{ stream_id, from_timestamp, to_timestamp }` — the stream to replay and the time window.

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

---

*End of guide. Memorize the short answers. Use the elaborations when they probe deeper. Good luck.*
