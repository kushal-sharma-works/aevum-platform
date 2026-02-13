import * as pulumi from "@pulumi/pulumi"
import { beforeEach, describe, expect, it } from "vitest"
import { createServicePolicies } from "../components/iam/policies"

describe("createServicePolicies", () => {
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

  it("does not use wildcard resource for event table policy", async () => {
    const policies = createServicePolicies({
      eventsTableArn: "arn:aws:dynamodb:eu-central-1:123456789012:table/aevum-dev-events",
      eventsTableStreamArn: "arn:aws:dynamodb:eu-central-1:123456789012:table/aevum-dev-events/stream/1",
      opensearchDomainArn: "arn:aws:es:eu-central-1:123456789012:domain/aevum-dev-search",
      documentDbClusterArn: "arn:aws:rds:eu-central-1:123456789012:cluster:aevum-dev-docdb",
      eventNotificationsQueueArn: "arn:aws:sqs:eu-central-1:123456789012:aevum-dev-event-notifications"
    })

    await new Promise<void>((resolve) => {
      pulumi.output(policies.eventTimelinePolicy.policy).apply((policyJson) => {
        const eventTimelinePolicyDoc = JSON.parse(policyJson)
        const resources = eventTimelinePolicyDoc.Statement[0].Resource as string[]
        expect(resources.some((resource) => resource === "*")).toBe(false)
        resolve()
      })
    })
  })
})
