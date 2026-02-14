import { DynamoDBStreamEvent, Context } from "aws-lambda"
import { SQSClient, SendMessageBatchCommand } from "@aws-sdk/client-sqs"

const sqs = new SQSClient({})
const QUEUE_URL = process.env.EVENT_NOTIFICATION_QUEUE_URL!

function chunk<T>(items: T[], size: number): T[][] {
  const batches: T[][] = []
  for (let i = 0; i < items.length; i += size) {
    batches.push(items.slice(i, i + size))
  }
  return batches
}

export async function handler(event: DynamoDBStreamEvent, _context: Context): Promise<void> {
  const messages = event.Records.filter((record) => record.eventName === "INSERT" || record.eventName === "MODIFY").map((record, index) => ({
    Id: `${record.eventID ?? "event"}-${index}`.slice(0, 80),
    MessageBody: JSON.stringify({
      eventName: record.eventName,
      eventSourceArn: record.eventSourceARN,
      keys: record.dynamodb?.Keys,
      newImage: record.dynamodb?.NewImage,
      oldImage: record.dynamodb?.OldImage,
      sequenceNumber: record.dynamodb?.SequenceNumber,
      approximateCreationDateTime: record.dynamodb?.ApproximateCreationDateTime
    })
  }))

  if (messages.length === 0) {
    console.log("No INSERT/MODIFY records to process")
    return
  }

  let sent = 0
  let failed = 0
  for (const batch of chunk(messages, 10)) {
    const response = await sqs.send(
      new SendMessageBatchCommand({
        QueueUrl: QUEUE_URL,
        Entries: batch
      })
    )
    sent += response.Successful?.length ?? 0
    failed += response.Failed?.length ?? 0
  }

  console.log(JSON.stringify({ totalMessages: messages.length, sent, failed }))
}
