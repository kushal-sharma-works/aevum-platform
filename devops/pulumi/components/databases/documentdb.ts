import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export interface DocumentDbOutputs {
  cluster: aws.docdb.Cluster
  connectionString: pulumi.Output<string>
  port: pulumi.Output<number | undefined>
}

export function createDocumentDb(
  config: InfraConfig,
  privateSubnetIds: pulumi.Input<string>[],
  documentDbSecurityGroupId: pulumi.Input<string>
): DocumentDbOutputs {
  const subnetGroup = new aws.docdb.SubnetGroup(resourceName("docdb-subnet-group"), {
    subnetIds: privateSubnetIds,
    tags: config.tags
  })

  const cluster = new aws.docdb.Cluster(resourceName("docdb"), {
    clusterIdentifier: resourceName("docdb"),
    engine: "docdb",
    engineVersion: "5.0.0",
    masterUsername: "aevum_admin",
    masterPassword: "ReplaceInSecretsManager123!",
    dbSubnetGroupName: subnetGroup.name,
    vpcSecurityGroupIds: [documentDbSecurityGroupId],
    backupRetentionPeriod: config.docdbBackupRetentionDays,
    storageEncrypted: true,
    skipFinalSnapshot: config.env !== "prod",
    applyImmediately: true,
    tags: {
      ...config.tags,
      Name: resourceName("docdb")
    }
  })

  for (let i = 0; i < config.docdbInstanceCount; i += 1) {
    new aws.docdb.ClusterInstance(resourceName(`docdb-instance-${i + 1}`), {
      clusterIdentifier: cluster.id,
      instanceClass: config.docdbInstanceType,
      applyImmediately: true,
      tags: {
        ...config.tags,
        Name: resourceName(`docdb-instance-${i + 1}`)
      }
    })
  }

  return {
    cluster,
    connectionString: pulumi.interpolate`mongodb://${cluster.masterUsername}:${cluster.masterPassword}@${cluster.endpoint}:27017/?tls=true&replicaSet=rs0&readPreference=secondaryPreferred&retryWrites=false`,
    port: cluster.port
  }
}
