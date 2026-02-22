# Deployment Guide

This repository is operated local-first. Cloud deployment assets exist under `devops/` as optional/manual paths.

## 1) Local deployment (default)

```bash
docker compose up -d --build
docker compose ps
```

Health checks:

```bash
curl -f http://localhost:8081/health
curl -f http://localhost:9091/admin/health
curl -f http://localhost:8080/health
```

Seed data:

```bash
./devops/scripts/seed-data.sh
```

## 2) Optional cloud deployment (manual)

### Prerequisites

- AWS account with IAM permissions for EKS, VPC, DynamoDB, OpenSearch, ECR, and related services
- Pulumi CLI
- AWS CLI
- `kubectl`
- Helm
- Docker

### Configure Pulumi Stack

```bash
cd devops/pulumi
pulumi stack init dev
pulumi config set aws:region eu-central-1
```

Set required stack configs (networking, cluster sizing, and service parameters) according to environment policy.

### Deploy Infrastructure

```bash
cd devops/pulumi
pulumi up
```

Validate outputs for EKS cluster, VPC, ECR repositories, and data services.

### Configure Kubernetes Context

```bash
aws eks update-kubeconfig --name <cluster-name> --region <region>
kubectl get nodes
```

### Build and Push Docker Images

```bash
docker build -t <ecr>/aevum-event-timeline:<tag> services/event-timeline
docker build -t <ecr>/aevum-decision-engine:<tag> services/decision-engine
docker build -t <ecr>/aevum-query-audit:<tag> services/query-audit
docker push <ecr>/aevum-event-timeline:<tag>
docker push <ecr>/aevum-decision-engine:<tag>
docker push <ecr>/aevum-query-audit:<tag>
```

### Deploy Services

### Option A: Helm

```bash
helm upgrade --install aevum devops/helm/umbrella -n aevum-sit --create-namespace
```

### Option B: ArgoCD

Apply ArgoCD project/application manifests and trigger sync.

### Verify Deployment

- `kubectl get pods -n aevum-sit`
- Check readiness/liveness probes
- Execute smoke requests against API endpoints

### Access Observability

- Port-forward or ingress to Grafana and Prometheus.
- Verify service dashboards and baseline telemetry.

### Troubleshooting

- **Pods pending**: check node capacity/taints and storage classes.
- **Image pull failures**: validate ECR auth and image tags.
- **DB connectivity failures**: check security groups, subnet routing, DNS.
- **Argo out-of-sync**: inspect app events and repo revision references.
