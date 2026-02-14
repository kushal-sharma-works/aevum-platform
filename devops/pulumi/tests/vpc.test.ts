import * as pulumi from "@pulumi/pulumi"
import { beforeEach, describe, expect, it } from "vitest"
import { createVpc } from "../components/networking/vpc"

describe("createVpc", () => {
  beforeEach(() => {
    pulumi.runtime.setMocks({
      newResource: function (args) {
        return {
          id: args.name,
          state: args.inputs
        }
      },
      call: function (args) {
        if (args.token === "aws:index/getAvailabilityZones:getAvailabilityZones") {
          return { names: ["eu-central-1a", "eu-central-1b", "eu-central-1c"] }
        }
        return args.inputs
      }
    })
  })

  it("creates three public and three private subnets", async () => {
    const vpc = createVpc({
      project: "aevum",
      env: "dev",
      awsRegion: "eu-central-1",
      vpcCidr: "10.0.0.0/16",
      natGatewayCount: 1,
      eksNodeCount: 2,
      eksMinSize: 1,
      eksMaxSize: 3,
      eksInstanceType: "t3.medium",
      opensearchInstanceType: "t3.small.search",
      opensearchInstanceCount: 1,
      opensearchVolumeSize: 10,
      docdbInstanceType: "db.t3.medium",
      docdbInstanceCount: 1,
      docdbBackupRetentionDays: 7,
      frontendDomainAliases: [],
      tags: { Project: "aevum", Environment: "dev", ManagedBy: "pulumi", Team: "platform" }
    })

    expect(vpc.publicSubnetIds).toHaveLength(3)
    expect(vpc.privateSubnetIds).toHaveLength(3)
  })
})
