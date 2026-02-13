import * as aws from "@pulumi/aws"
import * as k8s from "@pulumi/kubernetes"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createGrafanaMonitoring(
  config: InfraConfig,
  kubeProvider: k8s.Provider,
  oidcProviderArn: pulumi.Input<string>,
  oidcProviderUrl: pulumi.Input<string>
): { serviceAccount: k8s.core.v1.ServiceAccount; role: aws.iam.Role } {
  const role = new aws.iam.Role(resourceName("grafana-irsa-role"), {
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
                [`${issuer}:sub`]: "system:serviceaccount:monitoring:grafana",
                [`${issuer}:aud`]: "sts.amazonaws.com"
              }
            }
          }
        ]
      })
    }),
    tags: config.tags
  })

  const policy = new aws.iam.Policy(resourceName("grafana-cloudwatch-policy"), {
    policy: JSON.stringify({
      Version: "2012-10-17",
      Statement: [
        {
          Effect: "Allow",
          Action: [
            "cloudwatch:GetMetricData",
            "cloudwatch:ListMetrics",
            "logs:DescribeLogGroups",
            "logs:StartQuery",
            "logs:GetQueryResults"
          ],
          Resource: "*"
        }
      ]
    })
  })

  new aws.iam.RolePolicyAttachment(resourceName("grafana-policy-attach"), {
    role: role.name,
    policyArn: policy.arn
  })

  const serviceAccount = new k8s.core.v1.ServiceAccount(resourceName("grafana-sa"), {
    metadata: {
      name: "grafana",
      namespace: "monitoring",
      annotations: {
        "eks.amazonaws.com/role-arn": role.arn
      }
    }
  }, { provider: kubeProvider })

  return { serviceAccount, role }
}
