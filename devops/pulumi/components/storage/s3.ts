import * as aws from "@pulumi/aws"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createFrontendBucket(config: InfraConfig): aws.s3.BucketV2 {
  const bucket = new aws.s3.BucketV2(resourceName("frontend"), {
    bucket: resourceName("frontend"),
    tags: {
      ...config.tags,
      Name: resourceName("frontend")
    }
  })

  new aws.s3.BucketPublicAccessBlock(resourceName("frontend-pab"), {
    bucket: bucket.id,
    blockPublicAcls: true,
    blockPublicPolicy: true,
    ignorePublicAcls: true,
    restrictPublicBuckets: true
  })

  new aws.s3.BucketVersioningV2(resourceName("frontend-versioning"), {
    bucket: bucket.id,
    versioningConfiguration: {
      status: "Enabled"
    }
  })

  return bucket
}
