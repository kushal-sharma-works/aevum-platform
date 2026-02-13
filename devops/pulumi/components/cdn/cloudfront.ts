import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createFrontendCdn(
  config: InfraConfig,
  frontendBucket: aws.s3.BucketV2
): { distribution: aws.cloudfront.Distribution; oac: aws.cloudfront.OriginAccessControl } {
  const oac = new aws.cloudfront.OriginAccessControl(resourceName("frontend-oac"), {
    description: "OAC for Aevum frontend bucket",
    originAccessControlOriginType: "s3",
    signingBehavior: "always",
    signingProtocol: "sigv4"
  })

  const distribution = new aws.cloudfront.Distribution(resourceName("frontend-cdn"), {
    enabled: true,
    defaultRootObject: "index.html",
    aliases: config.frontendDomainAliases,
    origins: [
      {
        domainName: frontendBucket.bucketRegionalDomainName,
        originId: frontendBucket.arn,
        originAccessControlId: oac.id
      }
    ],
    defaultCacheBehavior: {
      targetOriginId: frontendBucket.arn,
      viewerProtocolPolicy: "redirect-to-https",
      allowedMethods: ["GET", "HEAD", "OPTIONS"],
      cachedMethods: ["GET", "HEAD", "OPTIONS"],
      compress: true,
      forwardedValues: {
        queryString: false,
        cookies: { forward: "none" }
      },
      minTtl: 0,
      defaultTtl: 3600,
      maxTtl: 31536000
    },
    orderedCacheBehaviors: [
      {
        pathPattern: "*.html",
        targetOriginId: frontendBucket.arn,
        viewerProtocolPolicy: "redirect-to-https",
        allowedMethods: ["GET", "HEAD", "OPTIONS"],
        cachedMethods: ["GET", "HEAD", "OPTIONS"],
        compress: true,
        forwardedValues: {
          queryString: false,
          cookies: { forward: "none" }
        },
        minTtl: 0,
        defaultTtl: 3600,
        maxTtl: 3600
      },
      {
        pathPattern: "assets/*",
        targetOriginId: frontendBucket.arn,
        viewerProtocolPolicy: "redirect-to-https",
        allowedMethods: ["GET", "HEAD", "OPTIONS"],
        cachedMethods: ["GET", "HEAD", "OPTIONS"],
        compress: true,
        forwardedValues: {
          queryString: false,
          cookies: { forward: "none" }
        },
        minTtl: 0,
        defaultTtl: 86400,
        maxTtl: 31536000
      }
    ],
    restrictions: {
      geoRestriction: {
        restrictionType: "none"
      }
    },
    customErrorResponses: [
      { errorCode: 403, responseCode: 200, responsePagePath: "/index.html", errorCachingMinTtl: 0 },
      { errorCode: 404, responseCode: 200, responsePagePath: "/index.html", errorCachingMinTtl: 0 }
    ],
    viewerCertificate: {
      cloudfrontDefaultCertificate: true
    },
    tags: {
      ...config.tags,
      Name: resourceName("frontend-cdn")
    }
  })

  const bucketPolicy = pulumi
    .all([frontendBucket.arn, distribution.arn])
    .apply(([bucketArn, distributionArn]) =>
      JSON.stringify({
        Version: "2012-10-17",
        Statement: [
          {
            Sid: "AllowCloudFrontServicePrincipalReadOnly",
            Effect: "Allow",
            Principal: {
              Service: "cloudfront.amazonaws.com"
            },
            Action: ["s3:GetObject"],
            Resource: [`${bucketArn}/*`],
            Condition: {
              StringEquals: {
                "AWS:SourceArn": distributionArn
              }
            }
          }
        ]
      })
    )

  new aws.s3.BucketPolicy(resourceName("frontend-bucket-policy"), {
    bucket: frontendBucket.id,
    policy: bucketPolicy
  })

  return { distribution, oac }
}
