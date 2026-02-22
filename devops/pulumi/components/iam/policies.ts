import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { resourceName } from "../../utils/naming"

export interface ServiceAccessPolicies {
  eventTimelinePolicy: aws.iam.Policy
  decisionEnginePolicy: aws.iam.Policy
  queryAuditPolicy: aws.iam.Policy
  lambdaDdbSqsPolicy: aws.iam.Policy
}

export function createServicePolicies(args: {
  eventsTableArn: pulumi.Input<string>
  eventsTableStreamArn: pulumi.Input<string>
  opensearchDomainArn: pulumi.Input<string>
  documentDbClusterArn: pulumi.Input<string>
  eventNotificationsQueueArn: pulumi.Input<string>
}): ServiceAccessPolicies {
  const eventTimelinePolicy = new aws.iam.Policy(resourceName("event-timeline-policy"), {
    policy: pulumi.all([args.eventsTableArn]).apply(([tableArn]) =>
      JSON.stringify({
        Version: "2012-10-17",
        Statement: [
          {
            Effect: "Allow",
            Action: ["dynamodb:GetItem", "dynamodb:PutItem", "dynamodb:UpdateItem", "dynamodb:Query", "dynamodb:Scan"],
            Resource: [tableArn, `${tableArn}/index/*`]
          }
        ]
      })
    )
  })

  const decisionEnginePolicy = new aws.iam.Policy(resourceName("decision-engine-policy"), {
    policy: pulumi.all([args.eventsTableArn, args.documentDbClusterArn]).apply(([tableArn, docdbArn]) =>
      JSON.stringify({
        Version: "2012-10-17",
        Statement: [
          {
            Effect: "Allow",
            Action: ["dynamodb:GetItem", "dynamodb:Query", "dynamodb:Scan"],
            Resource: [tableArn, `${tableArn}/index/*`]
          },
          {
            Effect: "Allow",
            Action: ["docdb-elastic:Connect"],
            Resource: [docdbArn]
          }
        ]
      })
    )
  })

  const queryAuditPolicy = new aws.iam.Policy(resourceName("query-audit-policy"), {
    policy: pulumi.all([args.eventsTableArn, args.opensearchDomainArn]).apply(([tableArn, searchArn]) =>
      JSON.stringify({
        Version: "2012-10-17",
        Statement: [
          {
            Effect: "Allow",
            Action: ["dynamodb:GetItem", "dynamodb:Query", "dynamodb:Scan"],
            Resource: [tableArn, `${tableArn}/index/*`]
          },
          {
            Effect: "Allow",
            Action: ["es:ESHttpGet", "es:ESHttpPost", "es:ESHttpPut", "es:ESHttpDelete"],
            Resource: [`${searchArn}/*`]
          }
        ]
      })
    )
  })

  const lambdaDdbSqsPolicy = new aws.iam.Policy(resourceName("lambda-ddb-sqs-policy"), {
    policy: pulumi.all([args.eventsTableStreamArn, args.eventNotificationsQueueArn]).apply(([streamArn, queueArn]) =>
      JSON.stringify({
        Version: "2012-10-17",
        Statement: [
          {
            Effect: "Allow",
            Action: [
              "dynamodb:DescribeStream",
              "dynamodb:GetRecords",
              "dynamodb:GetShardIterator",
              "dynamodb:ListStreams"
            ],
            Resource: [streamArn]
          },
          {
            Effect: "Allow",
            Action: ["sqs:SendMessage", "sqs:SendMessageBatch"],
            Resource: [queueArn]
          }
        ]
      })
    )
  })

  return {
    eventTimelinePolicy,
    decisionEnginePolicy,
    queryAuditPolicy,
    lambdaDdbSqsPolicy
  }
}
