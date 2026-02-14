import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export function createFrontendDns(
  config: InfraConfig,
  cloudfrontDomainName: pulumi.Input<string>
): { zone?: aws.route53.Zone; record?: aws.route53.Record } {
  if (!config.route53ZoneName || !config.route53RecordName) {
    return {}
  }

  const zone = new aws.route53.Zone(resourceName("hosted-zone"), {
    name: config.route53ZoneName,
    tags: config.tags
  })

  const record = new aws.route53.Record(resourceName("frontend-record"), {
    zoneId: zone.zoneId,
    name: config.route53RecordName,
    type: "CNAME",
    ttl: 60,
    records: [cloudfrontDomainName]
  })

  return { zone, record }
}
