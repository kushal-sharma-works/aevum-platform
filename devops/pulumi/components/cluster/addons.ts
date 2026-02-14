import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createEksAddons(config: InfraConfig, clusterName: pulumi.Input<string>): aws.eks.Addon[] {
  const addonNames = ["coredns", "kube-proxy", "vpc-cni", "aws-ebs-csi-driver"]
  return addonNames.map((addonName) =>
    new aws.eks.Addon(resourceName(`addon-${addonName}`), {
      addonName,
      clusterName,
      resolveConflictsOnCreate: "OVERWRITE",
      resolveConflictsOnUpdate: "OVERWRITE",
      tags: config.tags
    })
  )
}
