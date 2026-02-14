import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig } from "../../types"
import { resourceName } from "../../utils/naming"

export interface SecurityGroupOutputs {
  eksNodeSgId: pulumi.Output<string>
  opensearchSgId: pulumi.Output<string>
  documentDbSgId: pulumi.Output<string>
  albSgId: pulumi.Output<string>
}

export function createSecurityGroups(config: InfraConfig, vpcId: pulumi.Input<string>): SecurityGroupOutputs {
  const eksNodeSg = new aws.ec2.SecurityGroup(resourceName("eks-node-sg"), {
    vpcId,
    description: "EKS worker nodes",
    egress: [{ protocol: "-1", fromPort: 0, toPort: 0, cidrBlocks: ["0.0.0.0/0"] }],
    tags: { ...config.tags, Name: resourceName("eks-node-sg") }
  })

  const opensearchSg = new aws.ec2.SecurityGroup(resourceName("opensearch-sg"), {
    vpcId,
    description: "OpenSearch access from EKS nodes only",
    ingress: [{ protocol: "tcp", fromPort: 443, toPort: 443, securityGroups: [eksNodeSg.id] }],
    egress: [{ protocol: "-1", fromPort: 0, toPort: 0, cidrBlocks: ["0.0.0.0/0"] }],
    tags: { ...config.tags, Name: resourceName("opensearch-sg") }
  })

  const documentDbSg = new aws.ec2.SecurityGroup(resourceName("docdb-sg"), {
    vpcId,
    description: "DocumentDB access from EKS nodes only",
    ingress: [{ protocol: "tcp", fromPort: 27017, toPort: 27017, securityGroups: [eksNodeSg.id] }],
    egress: [{ protocol: "-1", fromPort: 0, toPort: 0, cidrBlocks: ["0.0.0.0/0"] }],
    tags: { ...config.tags, Name: resourceName("docdb-sg") }
  })

  const albSg = new aws.ec2.SecurityGroup(resourceName("alb-sg"), {
    vpcId,
    description: "Public ingress for platform endpoints",
    ingress: [
      { protocol: "tcp", fromPort: 80, toPort: 80, cidrBlocks: ["0.0.0.0/0"] },
      { protocol: "tcp", fromPort: 443, toPort: 443, cidrBlocks: ["0.0.0.0/0"] }
    ],
    egress: [{ protocol: "-1", fromPort: 0, toPort: 0, cidrBlocks: ["0.0.0.0/0"] }],
    tags: { ...config.tags, Name: resourceName("alb-sg") }
  })

  return {
    eksNodeSgId: eksNodeSg.id,
    opensearchSgId: opensearchSg.id,
    documentDbSgId: documentDbSg.id,
    albSgId: albSg.id
  }
}
