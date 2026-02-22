import * as aws from "@pulumi/aws"
import * as k8s from "@pulumi/kubernetes"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createPrometheusMonitoring(
  config: InfraConfig,
  kubeProvider: k8s.Provider,
  oidcProviderArn: pulumi.Input<string>,
  oidcProviderUrl: pulumi.Input<string>
): { namespace: k8s.core.v1.Namespace; serviceAccount: k8s.core.v1.ServiceAccount; role: aws.iam.Role } {
  const namespace = new k8s.core.v1.Namespace(resourceName("prometheus-ns"), {
    metadata: {
      name: "monitoring"
    }
  }, { provider: kubeProvider })

  const role = new aws.iam.Role(resourceName("prometheus-irsa-role"), {
    assumeRolePolicy: pulumi.all([oidcProviderArn, oidcProviderUrl]).apply(([arn, url]) => {
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
                [`${issuer}:sub`]: "system:serviceaccount:monitoring:prometheus",
                [`${issuer}:aud`]: "sts.amazonaws.com"
              }
            }
          }
        ]
      })
    }),
    tags: config.tags
  })

  const policy = new aws.iam.Policy(resourceName("prometheus-cloudwatch-policy"), {
    policy: JSON.stringify({
      Version: "2012-10-17",
      Statement: [
        {
          Effect: "Allow",
          Action: ["cloudwatch:GetMetricData", "cloudwatch:ListMetrics", "tag:GetResources"],
          Resource: "*"
        }
      ]
    })
  })

  new aws.iam.RolePolicyAttachment(resourceName("prometheus-policy-attach"), {
    role: role.name,
    policyArn: policy.arn
  })

  const serviceAccount = new k8s.core.v1.ServiceAccount(resourceName("prometheus-sa"), {
    metadata: {
      name: "prometheus",
      namespace: namespace.metadata.name,
      annotations: {
        "eks.amazonaws.com/role-arn": role.arn
      }
    }
  }, { provider: kubeProvider })

  return { namespace, serviceAccount, role }
}
