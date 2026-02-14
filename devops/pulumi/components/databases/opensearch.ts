import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createOpenSearchDomain(
  config: InfraConfig,
  privateSubnetIds: pulumi.Input<string>[],
  opensearchSecurityGroupId: pulumi.Input<string>
): aws.opensearch.Domain {
  return new aws.opensearch.Domain(resourceName("search"), {
    domainName: resourceName("search"),
    engineVersion: "OpenSearch_2.13",
    clusterConfig: {
      instanceType: config.opensearchInstanceType,
      instanceCount: config.opensearchInstanceCount,
      zoneAwarenessEnabled: config.env === "prod",
      zoneAwarenessConfig: config.env === "prod" ? { availabilityZoneCount: 3 } : undefined
    },
    ebsOptions: {
      ebsEnabled: true,
      volumeType: "gp3",
      volumeSize: config.opensearchVolumeSize
    },
    vpcOptions: {
      subnetIds: config.env === "prod" ? privateSubnetIds.slice(0, 3) : [privateSubnetIds[0]],
      securityGroupIds: [opensearchSecurityGroupId]
    },
    advancedSecurityOptions: {
      enabled: true,
      internalUserDatabaseEnabled: true,
      masterUserOptions: {
        masterUserName: "aevum-admin",
        masterUserPassword: "ReplaceInSecretsManager123!"
      }
    },
    encryptAtRest: { enabled: true },
    nodeToNodeEncryption: { enabled: true },
    domainEndpointOptions: {
      enforceHttps: true,
      tlsSecurityPolicy: "Policy-Min-TLS-1-2-2019-07"
    },
    accessPolicies: pulumi.output(aws.getCallerIdentity({})).apply((identity) =>
      JSON.stringify({
        Version: "2012-10-17",
        Statement: [
          {
            Sid: "AllowAccountAccess",
            Effect: "Allow",
            Principal: { AWS: `arn:aws:iam::${identity.accountId}:root` },
            Action: "es:ESHttp*",
            Resource: `arn:aws:es:${config.awsRegion}:${identity.accountId}:domain/${resourceName("search")}/*`
          }
        ]
      })
    ),
    tags: {
      ...config.tags,
      Name: resourceName("search")
    }
  })
}
