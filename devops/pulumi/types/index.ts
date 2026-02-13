import * as pulumi from "@pulumi/pulumi"

export type EnvName = "dev" | "staging" | "prod"

export interface InfraConfig {
  project: string
  env: EnvName
  awsRegion: string
  vpcCidr: string
  natGatewayCount: number
  eksNodeCount: number
  eksMinSize: number
  eksMaxSize: number
  eksInstanceType: string
  opensearchInstanceType: string
  opensearchInstanceCount: number
  opensearchVolumeSize: number
  docdbInstanceType: string
  docdbInstanceCount: number
  docdbBackupRetentionDays: number
  frontendDomainAliases: string[]
  route53ZoneName?: string
  route53RecordName?: string
  tags: Record<string, string>
}

export interface VpcOutputs {
  vpcId: pulumi.Output<string>
  publicSubnetIds: pulumi.Output<string>[]
  privateSubnetIds: pulumi.Output<string>[]
  flowLogGroupName: pulumi.Output<string>
}
