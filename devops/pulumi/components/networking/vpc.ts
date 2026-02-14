import * as aws from "@pulumi/aws"
import * as pulumi from "@pulumi/pulumi"
import { InfraConfig, VpcOutputs } from "../../types"
import { resourceName } from "../../utils/naming"

function subnetCidr(baseOctet: number, index: number): string {
  return `10.${baseOctet}.${index}.0/24`
}

export function createVpc(config: InfraConfig): VpcOutputs {
  const azs = aws.getAvailabilityZonesOutput({ state: "available" }).names
  const privateSubnets: aws.ec2.Subnet[] = []
  const publicSubnets: aws.ec2.Subnet[] = []

  const vpc = new aws.ec2.Vpc(resourceName("vpc"), {
    cidrBlock: config.vpcCidr,
    enableDnsHostnames: true,
    enableDnsSupport: true,
    tags: {
      ...config.tags,
      Name: resourceName("vpc")
    }
  })

  const igw = new aws.ec2.InternetGateway(resourceName("igw"), {
    vpcId: vpc.id,
    tags: {
      ...config.tags,
      Name: resourceName("igw")
    }
  })

  for (let i = 0; i < 3; i += 1) {
    const publicSubnet = new aws.ec2.Subnet(resourceName(`public-subnet-${i + 1}`), {
      vpcId: vpc.id,
      availabilityZone: azs.apply((z) => z[i]),
      mapPublicIpOnLaunch: true,
      cidrBlock: subnetCidr(0, i + 10),
      tags: {
        ...config.tags,
        Name: resourceName(`public-subnet-${i + 1}`),
        [`kubernetes.io/cluster/${resourceName("eks")}`]: "shared",
        "kubernetes.io/role/elb": "1"
      }
    })
    publicSubnets.push(publicSubnet)

    const privateSubnet = new aws.ec2.Subnet(resourceName(`private-subnet-${i + 1}`), {
      vpcId: vpc.id,
      availabilityZone: azs.apply((z) => z[i]),
      cidrBlock: subnetCidr(0, i + 20),
      tags: {
        ...config.tags,
        Name: resourceName(`private-subnet-${i + 1}`),
        [`kubernetes.io/cluster/${resourceName("eks")}`]: "shared",
        "kubernetes.io/role/internal-elb": "1"
      }
    })
    privateSubnets.push(privateSubnet)
  }

  const publicRouteTable = new aws.ec2.RouteTable(resourceName("public-rt"), {
    vpcId: vpc.id,
    routes: [{ cidrBlock: "0.0.0.0/0", gatewayId: igw.id }],
    tags: {
      ...config.tags,
      Name: resourceName("public-rt")
    }
  })

  publicSubnets.forEach((subnet, index) => {
    new aws.ec2.RouteTableAssociation(resourceName(`public-rta-${index + 1}`), {
      subnetId: subnet.id,
      routeTableId: publicRouteTable.id
    })
  })

  const natCount = Math.min(Math.max(config.natGatewayCount, 1), 3)
  const natGateways: aws.ec2.NatGateway[] = []

  for (let i = 0; i < natCount; i += 1) {
    const eip = new aws.ec2.Eip(resourceName(`nat-eip-${i + 1}`), {
      domain: "vpc",
      tags: {
        ...config.tags,
        Name: resourceName(`nat-eip-${i + 1}`)
      }
    })
    natGateways.push(
      new aws.ec2.NatGateway(resourceName(`nat-gw-${i + 1}`), {
        allocationId: eip.id,
        subnetId: publicSubnets[i].id,
        tags: {
          ...config.tags,
          Name: resourceName(`nat-gw-${i + 1}`)
        }
      })
    )
  }

  privateSubnets.forEach((subnet, index) => {
    const nat = natGateways[index % natGateways.length]
    const rt = new aws.ec2.RouteTable(resourceName(`private-rt-${index + 1}`), {
      vpcId: vpc.id,
      routes: [{ cidrBlock: "0.0.0.0/0", natGatewayId: nat.id }],
      tags: {
        ...config.tags,
        Name: resourceName(`private-rt-${index + 1}`)
      }
    })
    new aws.ec2.RouteTableAssociation(resourceName(`private-rta-${index + 1}`), {
      subnetId: subnet.id,
      routeTableId: rt.id
    })
  })

  const logGroup = new aws.cloudwatch.LogGroup(resourceName("vpc-flowlogs"), {
    retentionInDays: 30,
    tags: {
      ...config.tags,
      Name: resourceName("vpc-flowlogs")
    }
  })

  const flowLogRole = new aws.iam.Role(resourceName("vpc-flowlog-role"), {
    assumeRolePolicy: aws.iam.assumeRolePolicyForPrincipal({ Service: "vpc-flow-logs.amazonaws.com" }),
    tags: config.tags
  })

  new aws.iam.RolePolicyAttachment(resourceName("vpc-flowlog-role-policy"), {
    role: flowLogRole.name,
    policyArn: aws.iam.ManagedPolicies.CloudWatchLogsFullAccess
  })

  new aws.ec2.FlowLog(resourceName("vpc-flowlog"), {
    vpcId: vpc.id,
    iamRoleArn: flowLogRole.arn,
    logDestination: logGroup.arn,
    trafficType: "ALL",
    tags: config.tags
  })

  return {
    vpcId: vpc.id,
    publicSubnetIds: publicSubnets.map((s) => s.id),
    privateSubnetIds: privateSubnets.map((s) => s.id),
    flowLogGroupName: logGroup.name
  }
}
