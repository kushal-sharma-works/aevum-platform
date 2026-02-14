import * as pulumi from "@pulumi/pulumi"
import { beforeEach, describe, expect, it } from "vitest"
import { nodeGroupConfig } from "../components/cluster/node-group"

describe("nodeGroupConfig", () => {
  beforeEach(() => {
    pulumi.runtime.setMocks({
      newResource: function (args) {
        return { id: args.name, state: args.inputs }
      },
      call: function (args) {
        return args.inputs
      }
    })
  })

  it("uses stack-derived scaling values", () => {
    const nodeGroup = nodeGroupConfig({
      project: "aevum",
      env: "prod",
      awsRegion: "eu-central-1",
      vpcCidr: "10.2.0.0/16",
      natGatewayCount: 3,
      eksNodeCount: 3,
      eksMinSize: 3,
      eksMaxSize: 10,
      eksInstanceType: "m5.xlarge",
      opensearchInstanceType: "r6g.large.search",
      opensearchInstanceCount: 3,
      opensearchVolumeSize: 100,
      docdbInstanceType: "db.r6g.large",
      docdbInstanceCount: 3,
      docdbBackupRetentionDays: 35,
      frontendDomainAliases: [],
      tags: { Project: "aevum", Environment: "prod", ManagedBy: "pulumi", Team: "platform" }
    })

    expect(nodeGroup.desiredSize).toBe(3)
    expect(nodeGroup.maxSize).toBe(10)
    expect(nodeGroup.instanceType).toBe("m5.xlarge")
  })
})
