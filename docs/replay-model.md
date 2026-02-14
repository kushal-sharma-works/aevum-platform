# Replay Model

## What is Replay?

Replay reconstructs system behavior by re-processing historical events through the decision engine. In Aevum, replay is not an approximation process; it is a deterministic verification mechanism intended to produce equivalent outcomes for equivalent inputs.

## Determinism Contract

The following invariants are required:

1. **No ambient time**: all time values are injected explicitly. Go services use a Clock abstraction; .NET services use `TimeProvider`.
2. **Versioned rules**: evaluations bind to explicit rule versions, never implicit “latest”.
3. **Ordered events**: replay consumes events in strict sequence order within each stream.
4. **Idempotent output**: repeating the same event + rule version returns the existing decision.

## Replay Flow

1. Client submits replay request: `POST /admin/replay` with `{stream_id, from_timestamp, to_timestamp}`.
2. Event Timeline queries DynamoDB by stream and time range through sequence-aware index access.
3. Events are streamed through the replay engine via Go channels.
4. Replay engine invokes Decision Engine evaluate endpoint per event.
5. Decision Engine checks idempotency (event + rule identity) before evaluation.
6. Existing decision is returned when present; otherwise deterministic evaluation executes and persists.
7. Replay hash is compared with original decision hash (if historical decision exists).
8. Replay results are streamed progressively to clients (SSE).

## Hash Verification

Deterministic hash contract:

`SHA-256(EventId || RuleId || RuleVersion || EvaluatedAt.ToString("O"))`

A mismatch indicates either input drift (different rule version/timestamp), code-path drift, or an implementation defect.

## Simulation (“What-If”)

Replay can be executed against a candidate rule version instead of historical version to estimate behavioral impact before promoting rule changes. The resulting decision set can be diffed against the baseline for approval workflows.

## Limitations and Trade-offs

- Replay is sequential per stream to preserve deterministic ordering.
- Parallel replay across independent streams is safe.
- Cross-stream global ordering is intentionally not guaranteed.
- Runtime latency is non-deterministic; evaluation output should remain deterministic.
- Code changes in the Decision Engine may alter outputs despite stable historical data; this is a known and explicit trade-off.
