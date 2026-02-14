import * as pulumi from "@pulumi/pulumi"
import { InfraConfig, EnvName } from "./types"
import { defaultTags } from "./utils/tags"

function toEnvName(value: string): EnvName {
  if (value === "dev" || value === "staging" || value === "prod") {
    return value
  }
  throw new Error(`Unsupported stack name: ${value}. Expected dev|staging|prod.`)
}

function eksScaleForEnv(env: EnvName): { min: number; max: number; defaultCount: number; defaultType: string } {
  if (env === "dev") {
    return { min: 1, max: 3, defaultCount: 2, defaultType: "t3.medium" }
  }
  if (env === "staging") {
    return { min: 2, max: 5, defaultCount: 3, defaultType: "t3.large" }
  }
  return { min: 3, max: 10, defaultCount: 3, defaultType: "m5.xlarge" }
}

export function loadConfig(): InfraConfig {
  const stack = toEnvName(pulumi.getStack())
  const awsConfig = new pulumi.Config("aws")
  const config = new pulumi.Config("aevum")
  const eksScale = eksScaleForEnv(stack)

  return {
    project: "aevum",
    env: stack,
    awsRegion: awsConfig.require("region"),
    vpcCidr: config.get("vpcCidr") ?? (stack === "dev" ? "10.0.0.0/16" : stack === "staging" ? "10.1.0.0/16" : "10.2.0.0/16"),
    natGatewayCount: config.getNumber("natGatewayCount") ?? (stack === "prod" ? 3 : 1),
    eksNodeCount: config.getNumber("eksNodeCount") ?? eksScale.defaultCount,
    eksMinSize: eksScale.min,
    eksMaxSize: eksScale.max,
    eksInstanceType: config.get("eksInstanceType") ?? eksScale.defaultType,
    opensearchInstanceType: config.get("opensearchInstanceType") ?? (stack === "prod" ? "r6g.large.search" : "t3.small.search"),
    opensearchInstanceCount: config.getNumber("opensearchInstanceCount") ?? (stack === "prod" ? 3 : 1),
    opensearchVolumeSize: stack === "prod" ? 100 : 10,
    docdbInstanceType: config.get("docdbInstanceType") ?? (stack === "prod" ? "db.r6g.large" : "db.t3.medium"),
    docdbInstanceCount: config.getNumber("docdbInstanceCount") ?? (stack === "prod" ? 3 : 1),
    docdbBackupRetentionDays: stack === "prod" ? 35 : 7,
    frontendDomainAliases: config.getObject<string[]>("frontendDomainAliases") ?? [],
    route53ZoneName: config.get("route53ZoneName") ?? undefined,
    route53RecordName: config.get("route53RecordName") ?? undefined,
    tags: defaultTags()
  }
}
