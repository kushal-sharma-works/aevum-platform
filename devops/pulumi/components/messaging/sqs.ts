import * as aws from "@pulumi/aws"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export interface QueueOutputs {
  eventNotificationsQueue: aws.sqs.Queue
  eventNotificationsDlq: aws.sqs.Queue
  decisionRequestsQueue: aws.sqs.Queue
}

export function createQueues(config: InfraConfig): QueueOutputs {
  const dlq = new aws.sqs.Queue(resourceName("event-notifications-dlq"), {
    name: resourceName("event-notifications-dlq"),
    messageRetentionSeconds: 1209600,
    sqsManagedSseEnabled: true,
    tags: {
      ...config.tags,
      Name: resourceName("event-notifications-dlq")
    }
  })

  const eventQueue = new aws.sqs.Queue(resourceName("event-notifications"), {
    name: resourceName("event-notifications"),
    visibilityTimeoutSeconds: 60,
    messageRetentionSeconds: 1209600,
    sqsManagedSseEnabled: true,
    redrivePolicy: dlq.arn.apply((arn) =>
      JSON.stringify({
        deadLetterTargetArn: arn,
        maxReceiveCount: 3
      })
    ),
    tags: {
      ...config.tags,
      Name: resourceName("event-notifications")
    }
  })

  const decisionRequestsQueue = new aws.sqs.Queue(resourceName("decision-requests"), {
    name: resourceName("decision-requests"),
    visibilityTimeoutSeconds: 60,
    messageRetentionSeconds: 1209600,
    sqsManagedSseEnabled: true,
    tags: {
      ...config.tags,
      Name: resourceName("decision-requests")
    }
  })

  return {
    eventNotificationsQueue: eventQueue,
    eventNotificationsDlq: dlq,
    decisionRequestsQueue
  }
}
