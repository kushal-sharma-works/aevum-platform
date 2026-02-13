import * as eks from "@pulumi/eks"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"
import { nodeGroupConfig } from "./node-group"
import { createEksAddons } from "./addons"

export interface EksOutputs {
  clusterName: pulumi.Output<string>
  kubeconfig: pulumi.Output<string>
  oidcProviderArn: pulumi.Output<string>
  oidcProviderUrl: pulumi.Output<string>
}

export function createEksCluster(
  config: InfraConfig,
  privateSubnetIds: pulumi.Input<string>[],
  _nodeSecurityGroupId: pulumi.Input<string>
): EksOutputs {
  const ng = nodeGroupConfig(config)

  const cluster = new eks.Cluster(resourceName("eks"), {
    version: "1.31",
    vpcId: undefined,
    privateSubnetIds,
    endpointPrivateAccess: true,
    endpointPublicAccess: true,
    createOidcProvider: true,
    skipDefaultNodeGroup: true,
    enabledClusterLogTypes: ["api", "audit", "authenticator"],
    tags: {
      ...config.tags,
      Name: resourceName("eks")
    }
  })

  new eks.ManagedNodeGroup(resourceName("node-group"), {
    cluster,
    nodeRoleArn: cluster.instanceRoles.apply((roles) => roles[0].arn),
    subnetIds: privateSubnetIds,
    scalingConfig: {
      desiredSize: ng.desiredSize,
      minSize: ng.minSize,
      maxSize: ng.maxSize
    },
    instanceTypes: [ng.instanceType],
    labels: ng.labels
  })

  createEksAddons(config, cluster.eksCluster.name)

  return {
    clusterName: cluster.eksCluster.name,
    kubeconfig: cluster.kubeconfig.apply((cfg) => JSON.stringify(cfg)),
    oidcProviderArn: cluster.core.oidcProvider!.arn,
    oidcProviderUrl: cluster.core.oidcProvider!.url
  }
}
