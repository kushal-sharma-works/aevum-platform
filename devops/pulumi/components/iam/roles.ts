import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createEksNodeRole(config: InfraConfig): aws.iam.Role {
  const role = new aws.iam.Role(resourceName("eks-node-role"), {
    assumeRolePolicy: aws.iam.assumeRolePolicyForPrincipal({ Service: "ec2.amazonaws.com" }),
    tags: config.tags
  })

  ;[
    "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
    "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  ].forEach((policyArn, index) => {
    new aws.iam.RolePolicyAttachment(resourceName(`eks-node-role-attach-${index + 1}`), {
      role: role.name,
      policyArn
    })
  })

  return role
}

export function createLambdaExecutionRole(config: InfraConfig): aws.iam.Role {
  const role = new aws.iam.Role(resourceName("lambda-exec-role"), {
    assumeRolePolicy: aws.iam.assumeRolePolicyForPrincipal({ Service: "lambda.amazonaws.com" }),
    tags: config.tags
  })

  new aws.iam.RolePolicyAttachment(resourceName("lambda-basic-role-attach"), {
    role: role.name,
    policyArn: aws.iam.ManagedPolicies.AWSLambdaBasicExecutionRole
  })

  return role
}

function buildIrsaTrustPolicy(oidcProviderArn: pulumi.Input<string>, oidcProviderUrl: pulumi.Input<string>, namespace: string, serviceAccountName: string) {
  return pulumi.all([oidcProviderArn, oidcProviderUrl]).apply(([arn, url]) => {
    const issuer = url.replace("https://", "")
    return JSON.stringify({
      Version: "2012-10-17",
      Statement: [
        {
          Effect: "Allow",
          Principal: { Federated: arn },
          Action: "sts:AssumeRoleWithWebIdentity",
          Condition: {
            StringEquals: {
              [`${issuer}:sub`]: `system:serviceaccount:${namespace}:${serviceAccountName}`,
              [`${issuer}:aud`]: "sts.amazonaws.com"
            }
          }
        }
      ]
    })
  })
}

export function createIrsaServiceRoles(
  config: InfraConfig,
  oidcProviderArn: pulumi.Input<string>,
  oidcProviderUrl: pulumi.Input<string>
): Record<string, aws.iam.Role> {
  const services = [
    { key: "eventTimeline", roleName: "event-timeline-role", namespace: "aevum", sa: "event-timeline" },
    { key: "decisionEngine", roleName: "decision-engine-role", namespace: "aevum", sa: "decision-engine" },
    { key: "queryAudit", roleName: "query-audit-role", namespace: "aevum", sa: "query-audit" }
  ]

  return services.reduce<Record<string, aws.iam.Role>>((acc, svc) => {
    acc[svc.key] = new aws.iam.Role(resourceName(svc.roleName), {
      assumeRolePolicy: buildIrsaTrustPolicy(oidcProviderArn, oidcProviderUrl, svc.namespace, svc.sa),
      tags: config.tags
    })
    return acc
  }, {})
}
