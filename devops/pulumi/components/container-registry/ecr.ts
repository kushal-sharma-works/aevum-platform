import * as aws from "@pulumi/aws"
import { InfraConfig } from "../../types"

const repositories = [
  "aevum-event-timeline",
  "aevum-decision-engine",
  "aevum-query-audit",
  "aevum-ui"
]

export function createEcrRepositories(config: InfraConfig): Record<string, aws.ecr.Repository> {
  const lifecyclePolicyDocument = JSON.stringify({
    rules: [
      {
        rulePriority: 1,
        description: "Keep last 10 tagged images",
        selection: {
          tagStatus: "tagged",
          tagPrefixList: ["v", "release", "main", "dev", "staging", "prod"],
          countType: "imageCountMoreThan",
          countNumber: 10
        },
        action: { type: "expire" }
      },
      {
        rulePriority: 2,
        description: "Expire untagged images older than 7 days",
        selection: {
          tagStatus: "untagged",
          countType: "sinceImagePushed",
          countUnit: "days",
          countNumber: 7
        },
        action: { type: "expire" }
      }
    ]
  })

  return repositories.reduce<Record<string, aws.ecr.Repository>>((acc, repoName) => {
    const repo = new aws.ecr.Repository(repoName, {
      name: repoName,
      imageTagMutability: "MUTABLE",
      imageScanningConfiguration: { scanOnPush: true },
      encryptionConfigurations: [{ encryptionType: "AES256" }],
      tags: {
        ...config.tags,
        Name: repoName
      }
    })
    new aws.ecr.LifecyclePolicy(`${repoName}-lifecycle`, {
      repository: repo.name,
      policy: lifecyclePolicyDocument
    })
    acc[repoName] = repo
    return acc
  }, {})
}
