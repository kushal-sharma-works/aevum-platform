# Deployment Guide (AWS)

## 1) Prerequisites

- AWS account with IAM permissions for EKS, VPC, DynamoDB, OpenSearch, ECR, and related services
- Pulumi CLI
- AWS CLI
- `kubectl`
- Helm
- Docker

## 2) Configure Pulumi Stack

```bash
cd devops/pulumi
pulumi stack init dev
pulumi config set aws:region eu-central-1
```

Set required stack configs (networking, cluster sizing, and service parameters) according to environment policy.

## 3) Deploy Infrastructure

```bash
cd devops/pulumi
pulumi up
```

Validate outputs for EKS cluster, VPC, ECR repositories, and data services.

## 4) Configure Kubernetes Context

```bash
aws eks update-kubeconfig --name <cluster-name> --region <region>
kubectl get nodes
```

## 5) Build and Push Docker Images

```bash
docker build -t <ecr>/aevum-event-timeline:<tag> services/event-timeline
docker build -t <ecr>/aevum-decision-engine:<tag> services/decision-engine
docker build -t <ecr>/aevum-query-audit:<tag> services/query-audit
docker push <ecr>/aevum-event-timeline:<tag>
docker push <ecr>/aevum-decision-engine:<tag>
docker push <ecr>/aevum-query-audit:<tag>
```

## 6) Deploy Services

### Option A: Helm

```bash
helm upgrade --install aevum devops/helm/umbrella -n aevum-sit --create-namespace
```

### Option B: ArgoCD

Apply ArgoCD project/application manifests and trigger sync.

## 7) Verify Deployment

- `kubectl get pods -n aevum-sit`
- Check readiness/liveness probes
- Execute smoke requests against API endpoints

## 8) Access Observability

- Port-forward or ingress to Grafana and Prometheus.
- Verify service dashboards and baseline telemetry.

## 9) Troubleshooting

- **Pods pending**: check node capacity/taints and storage classes.
- **Image pull failures**: validate ECR auth and image tags.
- **DB connectivity failures**: check security groups, subnet routing, DNS.
- **Argo out-of-sync**: inspect app events and repo revision references.
