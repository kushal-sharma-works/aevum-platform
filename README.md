# Aevum

Aevum is a distributed, deterministic, time-aware decision platform.

The system is designed around immutable event ingestion, deterministic decision evaluation, and historical replay for auditability and analysis.

## Repository Structure

- `services/` – Backend microservices (Go, .NET)
- `frontend/` – Web UI (Vue + TypeScript)
	- `frontend/aevum-ui/` – Vue 3 + TypeScript + Pinia + Vite frontend
- `devops/` – Infrastructure and platform operations
	- `devops/pulumi/` – Pulumi TypeScript IaC for AWS
	- `devops/helm/` – Helm charts for all services + umbrella chart
	- `devops/argocd/` – ArgoCD project, applications, and applicationset
	- `devops/monitoring/` – OTel, Prometheus, and Grafana provisioning
	- `devops/k8s/` – Base namespaces, quotas, and limit ranges
	- `devops/scripts/` – Local setup, seed, and cluster port-forward helpers

## Local Development Workflow

Use the root Makefile and Docker Compose stack for one-command local setup.

```bash
make dev
```

This starts:

- Event Timeline: `http://localhost:8080` (admin mapped to `9091`)
- Decision Engine: `http://localhost:8081`
- Query & Audit: `http://localhost:8082`
- Frontend: `http://localhost:3000`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3001` (`admin/admin`)

Useful commands:

```bash
make dev-down
make test
make lint
make build
make seed
```

Optional helper scripts:

```bash
bash devops/scripts/local-setup.sh
bash devops/scripts/seed-data.sh
bash devops/scripts/port-forward.sh aevum-dev
```

## CI/CD Strategy

GitHub Actions workflows are in `.github/workflows/`:

- `ci.yml` – change-aware lint/test/build workflow
- `pr-check.yml` – Helm + Argo manifest validation and rendering checks
- `cd-dev.yml` – updates `values-dev.yaml` image tags on successful main CI
- `cd-staging.yml` – updates `values-staging.yaml` image tags on release
- `cd-prod.yml` – updates `values-prod.yaml` image tags on `v*.*.*` tags
- `infra.yml` – Pulumi preview on PR and Pulumi up on main

GitLab parity pipeline exists in `.gitlab-ci.yml` with equivalent stages for lint, test, build, and deploy.

ArgoCD pulls desired state from repo manifests under `devops/argocd/` and deploys Helm charts from `devops/helm/charts/`.

## Monitoring Setup

- OpenTelemetry Collector config: `devops/monitoring/otel-collector-config.yaml`
- Prometheus local scrape + rules: `devops/monitoring/prometheus/`
- Grafana datasource/provider + dashboards: `devops/monitoring/grafana/`

Dashboards include:

- Aevum overview
- Event Timeline service metrics
- Decision Engine service metrics
- Query & Audit service metrics
- Kubernetes infrastructure metrics

