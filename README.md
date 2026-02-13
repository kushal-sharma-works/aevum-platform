# Aevum

Aevum is a distributed, deterministic, time-aware decision platform.

The system is designed around immutable event ingestion, deterministic decision evaluation, and historical replay for auditability and analysis.

## Repository Structure

- `services/`   – Backend microservices (Go, .NET)
- `frontend/`   – Web UI (Vue + TypeScript)
	- `frontend/aevum-ui/` – Full PR4 application (Vue 3 + TypeScript + Pinia + Vite)
- `devops/`     – Infrastructure
	- `devops/pulumi/` – Full Pulumi TypeScript IaC for AWS (VPC, EKS, DynamoDB, OpenSearch, DocumentDB, ECR, Lambda, S3/CloudFront, IAM, monitoring)

