# Decision Engine Service

A deterministic rule evaluation engine built with C# .NET 9, designed for idempotent and traceable decision-making workflows.

## Features

- **Deterministic Evaluation**: SHA256-based hashing ensures identical inputs always produce the same decisions
- **Flexible Rule Engine**: Support for complex conditions with multiple operators (equals, greater than, contains, regex, etc.)
- **Versioned Rules**: Immutable rule versioning with full history tracking
- **Idempotency**: Hash-based deduplication prevents duplicate decisions
- **MongoDB Persistence**: Scalable storage for rules and decisions with optimized indexing
- **Event Timeline Integration**: Optional integration with Event Timeline service via Polly-protected HTTP client
- **Observability**: OpenTelemetry tracing/metrics, Prometheus metrics, Serilog structured logging
- **Production-Ready**: Comprehensive test suite, health checks, graceful shutdown, RFC 9457 Problem Details

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     API Layer                            │
│  (Minimal API, Endpoints, Middleware, Filters)           │
└─────────────────┬───────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────┐
│                Application Layer                         │
│  • DeterministicEvaluator (Core Logic)                   │
│  • EvaluationService                                     │
│  • RuleManagementService                                 │
│  • DTOs, Validators, Mappers                             │
└─────────────────┬───────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────┐
│              Infrastructure Layer                        │
│  • MongoDbRuleRepository                                 │
│  • MongoDbDecisionRepository                             │
│  • EventTimelineClient (Polly retry/circuit breaker)     │
└─────────────────┬───────────────────────────────────────┘
                  │
┌─────────────────▼───────────────────────────────────────┐
│                  Domain Layer                            │
│  • Immutable Records (Rule, Decision, EvaluationContext) │
│  • Interfaces, Enums, Exceptions                         │
└──────────────────────────────────────────────────────────┘
```

## Technology Stack

- **.NET 9**: Target framework with C# preview features (collection expressions, primary constructors)
- **MongoDB 3.2**: Document storage with compound indexes
- **FluentValidation 11.11**: Request validation
- **Serilog**: Structured logging
- **OpenTelemetry**: Distributed tracing and metrics
- **Prometheus**: Metrics exposition
- **Polly**: Resilience policies (retry, circuit breaker)
- **xUnit + FluentAssertions**: Testing framework
- **Testcontainers**: Integration testing with real MongoDB

## API Endpoints

### Rules Management

```
POST   /api/v1/rules              - Create new rule
GET    /api/v1/rules/{id}         - Get rule by ID (optionally versioned)
PUT    /api/v1/rules/{id}         - Update rule (creates new version)
DELETE /api/v1/rules/{id}         - Delete rule
GET    /api/v1/rules?status=...   - List rules by status
POST   /api/v1/rules/{id}/activate   - Activate rule
POST   /api/v1/rules/{id}/deactivate - Deactivate rule
```

### Decision Evaluation

```
POST   /api/v1/decisions/evaluate           - Evaluate decision against rule
GET    /api/v1/decisions/{id}               - Get decision by ID
GET    /api/v1/decisions/request/{reqId}    - Get decision by request ID (idempotency)
GET    /api/v1/decisions/rule/{ruleId}      - List decisions by rule
```

### System

```
GET    /health        - Health check
GET    /health/ready  - Readiness probe
GET    /health/live   - Liveness probe
GET    /metrics       - Prometheus metrics
```

## Configuration

Environment variables or `appsettings.json`:

```json
{
  "MongoDB": {
    "ConnectionString": "mongodb://localhost:27017",
    "DatabaseName": "decision_engine"
  },
  "EventTimeline": {
    "BaseUrl": "http://localhost:8080"
  },
  "Observability": {
    "ServiceName": "decision-engine",
    "ServiceVersion": "1.0.0",
    "OtlpEndpoint": "http://localhost:4317"
  }
}
```

## Quick Start

### Prerequisites

- .NET 9 SDK
- MongoDB 8.0+
- Docker (optional)

### Local Development

```bash
# Restore dependencies
make restore

# Build
make build

# Run tests
make test

# Run locally
make run

# Or with hot reload
make watch
```

### Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

## Example Usage

### Create a Rule

```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "High Value Transaction Rule",
    "description": "Approve transactions over $1000",
    "conditions": [{
      "field": "amount",
      "operator": "GreaterThan",
      "value": 1000
    }],
    "actions": [{
      "type": "StoreDecision",
      "parameters": { "action": "approve" },
      "order": 1
    }],
    "priority": 10
  }'
```

### Activate Rule

```bash
curl -X POST http://localhost:8080/api/v1/rules/{ruleId}/activate
```

### Evaluate Decision

```bash
curl -X POST http://localhost:8080/api/v1/decisions/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "ruleId": "rule-123",
    "context": {
      "amount": 1500,
      "currency": "USD"
    },
    "requestId": "req-456"
  }'
```

## Testing

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests only (requires Docker)
make test-integration

# With coverage
make test-coverage
```

## Deterministic Evaluation

The DeterministicEvaluator guarantees:

1. **Hash Consistency**: Same rule + context → same SHA256 hash
2. **Idempotency**: Duplicate requests return cached decisions
3. **Immutability**: All domain models are immutable records
4. **Time Independence**: Uses injected TimeProvider (no DateTime.Now)

## Rule Operators

- **Comparison**: Equals, NotEquals, GreaterThan, GreaterThanOrEqual, LessThan, LessThanOrEqual
- **String**: Contains, NotContains, StartsWith, EndsWith
- **Collection**: In, NotIn
- **Pattern**: Regex

## Observability

- **Metrics**: Prometheus metrics at `/metrics`
- **Tracing**: OpenTelemetry OTLP export
- **Logging**: Serilog structured logs (JSON console output)
- **Health checks**: Kubernetes-ready probes

## License

Proprietary - Aevum Platform
