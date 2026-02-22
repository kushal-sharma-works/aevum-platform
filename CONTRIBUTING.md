# Contributing to Aevum

Thank you for contributing to Aevum. This guide defines the engineering standards used across backend services, frontend applications, and infrastructure.

## Development Environment

1. Install prerequisites:
   - Docker + Docker Compose
   - Go 1.22+
   - .NET 9 SDK
   - Node.js 22+
2. Clone repository and install dependencies for the area you are modifying.
3. Start local stack:

```bash
docker compose up -d --build
```

4. Validate changes before opening a PR:

```bash
# per-service examples
cd services/event-timeline && make test-unit
cd services/decision-engine && make test-unit
cd services/query-audit && go test ./...
cd frontend/aevum-ui && npm run lint && npm run test
```

## Branch Naming Convention

Use one of the following prefixes:

- `feat/` for new features
- `fix/` for bug fixes
- `docs/` for documentation-only changes
- `refactor/` for internal code restructuring
- `chore/` for maintenance or tooling work

Example:

```text
feat/add-replay-impact-analysis
```

## Commit Message Convention

Use Conventional Commits:

- `feat:`
- `fix:`
- `docs:`
- `refactor:`
- `test:`
- `chore:`

Example:

```text
feat: add deterministic replay hash verification endpoint
```

## Pull Request Process

1. Branch from `main`.
2. Keep PR scope focused and reviewable.
3. Ensure CI is green before requesting review.
4. Add or update tests for behavior changes.
5. Update related documentation (API specs, architecture, ADRs) when design changes.
6. Request review from service owners.

## Code Style and Quality

- **Go**: enforce `golangci-lint`, keep handlers thin and business logic in services.
- **C#/.NET**: enforce `dotnet format` and analyzer warnings; prefer immutable models and explicit validation.
- **TypeScript/Vue**: enforce ESLint + Prettier; maintain strict typing and composable components.

## Testing Expectations

- Unit tests are required for business logic changes.
- Integration tests are required for persistence and HTTP contract changes.
- Smoke tests should verify startup and health/readiness paths.
- Do not merge PRs that lower coverage in core modules without explicit rationale.

## Documentation Expectations

Documentation is part of the deliverable. When changing architecture, APIs, runbooks, or workflows, update the relevant files under `docs/` in the same PR.
