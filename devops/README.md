# DevOps Runbook

This repository now includes a complete deployment/testing suite for:

- Local development and smoke testing via Docker Compose
- SIT deployment via Kubernetes manifests (Kustomize)
- GitOps deployment via ArgoCD (Application + ApplicationSet)

## Local: one-command stack

From repository root:

```bash
./scripts/local-setup.sh
```

What starts:

- MongoDB (`localhost:27017`)
- DynamoDB Local (`localhost:8000`)
- Event Timeline (`localhost:8081`, admin `localhost:9091`)
- Decision Engine (`localhost:8080`)

Quick smoke checks:

```bash
curl -f http://localhost:8080/health
curl -f http://localhost:9091/admin/health
```

## Local tests (all services)

Unit-only (default):

```bash
./scripts/local-test.sh
```

Unit + integration (opt-in):

```bash
./scripts/local-test.sh --integration
```

## SIT deploy (kubectl)

Prereqs:

- Access to target Kubernetes cluster
- `kubectl` configured to SIT context

Deploy:

```bash
./scripts/sit-deploy.sh
```

Manifests path:

- Base: `devops/k8s/base`
- SIT overlay: `devops/k8s/overlays/sit`

## ArgoCD

Apply project + apps:

```bash
kubectl apply -f devops/argocd/project.yaml
kubectl apply -f devops/argocd/apps/aevum-sit.yaml
kubectl apply -f devops/argocd/applicationset.yaml
```

## Image publishing notes

SIT manifests reference:

- `ghcr.io/kushal-sharma-works/aevum-decision-engine:latest`
- `ghcr.io/kushal-sharma-works/aevum-event-timeline:latest`

Publish these images before deploying to SIT, or patch image names/tags in `devops/k8s/overlays/sit/kustomization.yaml`.
