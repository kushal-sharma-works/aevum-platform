import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export interface FanoutLambdaOutputs {
  lambda: aws.lambda.Function
}

export function createEventFanoutLambda(args: {
  config: InfraConfig
  lambdaRoleArn: pulumi.Input<string>
  tableStreamArn: pulumi.Input<string>
  eventQueueArn: pulumi.Input<string>
  eventQueueUrl: pulumi.Input<string>
  deadLetterQueueArn: pulumi.Input<string>
}): FanoutLambdaOutputs {
  const fn = new aws.lambda.Function(resourceName("event-fanout"), {
    name: resourceName("event-fanout"),
    role: args.lambdaRoleArn,
    runtime: "nodejs20.x",
    handler: "index.handler",
    memorySize: 256,
    timeout: 30,
    code: new pulumi.asset.AssetArchive({
      "index.ts": new pulumi.asset.FileAsset("./components/compute/lambda-code/fanout/index.ts"),
      "package.json": new pulumi.asset.FileAsset("./components/compute/lambda-code/fanout/package.json")
    }),
    environment: {
      variables: {
        EVENT_NOTIFICATION_QUEUE_URL: args.eventQueueUrl
      }
    },
    deadLetterConfig: {
      targetArn: args.deadLetterQueueArn
    },
    tags: {
      ...args.config.tags,
      Name: resourceName("event-fanout")
    }
  })

  new aws.lambda.EventSourceMapping(resourceName("event-fanout-stream-mapping"), {
    eventSourceArn: args.tableStreamArn,
    functionName: fn.arn,
    batchSize: 100,
    maximumBatchingWindowInSeconds: 5,
    startingPosition: "LATEST"
  })

  new aws.lambda.Permission(resourceName("event-fanout-stream-permission"), {
    action: "lambda:InvokeFunction",
    function: fn.name,
    principal: "dynamodb.amazonaws.com",
    sourceArn: args.tableStreamArn
  })

  return { lambda: fn }
}
