# DevOps Runbook

This repository includes deployment/testing assets for:

- Local development and smoke testing via Docker Compose
- optional SIT deployment via Kubernetes manifests (Kustomize)
- optional GitOps deployment via ArgoCD (Application + ApplicationSet)

## Local: one-command stack

From repository root:

```bash
./devops/scripts/local-setup.sh
```

What starts:

- The setup script prepares local DynamoDB table(s).
- Start services with root compose: `docker compose up -d --build`

Quick smoke checks:

```bash
curl -f http://localhost:8080/health
curl -f http://localhost:9091/admin/health
```

## Local tests (all services)

Unit-only (default):

Run service-level tests from each service directory (`make test-unit` / language-native test command).

Unit + integration (opt-in):

Integration tests are optional and service-specific.

## SIT deploy (kubectl)

Prereqs:

- Access to target Kubernetes cluster
- `kubectl` configured to SIT context

Deploy (manual):

```bash
kubectl apply -k devops/k8s/overlays/sit
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
