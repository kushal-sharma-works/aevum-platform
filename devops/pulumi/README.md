# Aevum Pulumi Infrastructure

This directory contains complete AWS infrastructure as code for Aevum using Pulumi TypeScript.

## What it provisions

- Networking: VPC, public/private subnets, NAT, flow logs, security groups
- Compute: EKS cluster, managed node group, addons, Lambda event fanout
- Data: DynamoDB, OpenSearch, DocumentDB
- Messaging: SQS queues and DLQ
- Delivery: ECR repos, S3 frontend bucket, CloudFront with OAC
- Observability: Prometheus and Grafana namespaces + IRSA roles
- Optional DNS: Route53 zone and frontend record

## Stacks

- `dev` (`Pulumi.dev.yaml`)
- `staging` (`Pulumi.staging.yaml`)
- `prod` (`Pulumi.prod.yaml`)

## Usage

1. Install dependencies:

   `npm install`

2. Select stack:

   `pulumi stack select dev`

3. Preview/apply:

   `make preview`

   `make up`

4. Run tests:

   `make test`

## Important

- Replace placeholder credentials for OpenSearch and DocumentDB with secrets via Pulumi config and/or AWS Secrets Manager before deploy.
- All resources include standard tags from `utils/tags.ts`.
- Critical outputs are exported from `index.ts` for CI/CD and Helm integrations.
