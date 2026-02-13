import * as aws from "@pulumi/aws"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createEventsTable(config: InfraConfig): aws.dynamodb.Table {
  return new aws.dynamodb.Table(resourceName("events"), {
    name: resourceName("events"),
    billingMode: "PAY_PER_REQUEST",
    hashKey: "PK",
    rangeKey: "SK",
    attributes: [
      { name: "PK", type: "S" },
      { name: "SK", type: "S" },
      { name: "GSI1PK", type: "S" },
      { name: "GSI1SK", type: "N" },
      { name: "GSI2PK", type: "S" }
    ],
    globalSecondaryIndexes: [
      {
        name: "stream-sequence-index",
        hashKey: "GSI1PK",
        rangeKey: "GSI1SK",
        projectionType: "ALL"
      },
      {
        name: "idempotency-index",
        hashKey: "GSI2PK",
        projectionType: "KEYS_ONLY"
      }
    ],
    streamEnabled: true,
    streamViewType: "NEW_AND_OLD_IMAGES",
    pointInTimeRecovery: { enabled: true },
    serverSideEncryption: { enabled: true },
    tags: {
      ...config.tags,
      Name: resourceName("events")
    }
  })
}
