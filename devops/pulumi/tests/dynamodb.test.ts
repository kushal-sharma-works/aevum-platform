import * as pulumi from "@pulumi/pulumi"
import * as fs from "node:fs"
import * as path from "node:path"
import { beforeEach, describe, expect, it } from "vitest"
import { createEventsTable } from "../components/databases/dynamodb"

describe("createEventsTable", () => {
  beforeEach(() => {
    pulumi.runtime.setMocks({
      newResource: function (args) {
        return {
          id: args.name,
          state: args.inputs
        }
      },
      call: function (args) {
        return args.inputs
      }
    })
  })

  it("configures required GSIs", async () => {
    const table = createEventsTable({
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

    await new Promise<void>((resolve) => {
      pulumi.output(table.name).apply((tableName) => {
        expect(tableName).toMatch(/^aevum-.+-events$/)
        resolve()
      })
    })

    const source = fs.readFileSync(path.resolve(__dirname, "../components/databases/dynamodb.ts"), "utf8")
    expect(source).toContain("stream-sequence-index")
    expect(source).toContain("idempotency-index")
  })
})
