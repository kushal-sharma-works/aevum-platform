import * as aws from "@pulumi/aws"
import * as k8s from "@pulumi/kubernetes"
import { loadConfig } from "./config"
import { createEksCluster } from "./components/cluster/eks"
import { createEcrRepositories } from "./components/container-registry/ecr"
import { createEventFanoutLambda } from "./components/compute/lambda"
import { createFrontendCdn } from "./components/cdn/cloudfront"
import { createDocumentDb } from "./components/databases/documentdb"
import { createEventsTable } from "./components/databases/dynamodb"
import { createOpenSearchDomain } from "./components/databases/opensearch"
import { createFrontendDns } from "./components/dns/route53"
import { createServicePolicies } from "./components/iam/policies"
import { createEksNodeRole, createIrsaServiceRoles, createLambdaExecutionRole } from "./components/iam/roles"
import { createQueues } from "./components/messaging/sqs"
import { createGrafanaMonitoring } from "./components/monitoring/grafana"
import { createPrometheusMonitoring } from "./components/monitoring/prometheus"
import { createSecurityGroups } from "./components/networking/security-groups"
import { createVpc } from "./components/networking/vpc"
import { createFrontendBucket } from "./components/storage/s3"

const config = loadConfig()

const vpc = createVpc(config)
const securityGroups = createSecurityGroups(config, vpc.vpcId)

const eventsTable = createEventsTable(config)
const opensearch = createOpenSearchDomain(config, vpc.privateSubnetIds, securityGroups.opensearchSgId)
const docdb = createDocumentDb(config, vpc.privateSubnetIds, securityGroups.documentDbSgId)

const ecrRepositories = createEcrRepositories(config)

const eksNodeRole = createEksNodeRole(config)
const cluster = createEksCluster(config, vpc.privateSubnetIds, securityGroups.eksNodeSgId)

const queues = createQueues(config)
const lambdaRole = createLambdaExecutionRole(config)

const policies = createServicePolicies({
  eventsTableArn: eventsTable.arn,
  eventsTableStreamArn: eventsTable.streamArn,
  opensearchDomainArn: opensearch.arn,
  documentDbClusterArn: docdb.cluster.arn,
  eventNotificationsQueueArn: queues.eventNotificationsQueue.arn
})

new aws.iam.RolePolicyAttachment("lambda-ddb-sqs-policy-attach", {
  role: lambdaRole.name,
  policyArn: policies.lambdaDdbSqsPolicy.arn
})

const irsaRoles = createIrsaServiceRoles(config, cluster.oidcProviderArn, cluster.oidcProviderUrl)
new aws.iam.RolePolicyAttachment("event-timeline-policy-attach", {
  role: irsaRoles.eventTimeline.name,
  policyArn: policies.eventTimelinePolicy.arn
})
new aws.iam.RolePolicyAttachment("decision-engine-policy-attach", {
  role: irsaRoles.decisionEngine.name,
  policyArn: policies.decisionEnginePolicy.arn
})
new aws.iam.RolePolicyAttachment("query-audit-policy-attach", {
  role: irsaRoles.queryAudit.name,
  policyArn: policies.queryAuditPolicy.arn
})

const fanoutLambda = createEventFanoutLambda({
  config,
  lambdaRoleArn: lambdaRole.arn,
  tableStreamArn: eventsTable.streamArn,
  eventQueueArn: queues.eventNotificationsQueue.arn,
  eventQueueUrl: queues.eventNotificationsQueue.url,
  deadLetterQueueArn: queues.eventNotificationsDlq.arn
})

const frontendBucket = createFrontendBucket(config)
const cloudfront = createFrontendCdn(config, frontendBucket)
const dns = createFrontendDns(config, cloudfront.distribution.domainName)

const k8sProvider = new k8s.Provider("aevum-k8s", { kubeconfig: cluster.kubeconfig })
const prometheus = createPrometheusMonitoring(config, k8sProvider, cluster.oidcProviderArn, cluster.oidcProviderUrl)
const grafana = createGrafanaMonitoring(config, k8sProvider, cluster.oidcProviderArn, cluster.oidcProviderUrl)

export const vpcId = vpc.vpcId
export const publicSubnetIds = vpc.publicSubnetIds
export const privateSubnetIds = vpc.privateSubnetIds
export const flowLogGroupName = vpc.flowLogGroupName
export const eksClusterName = cluster.clusterName
export const kubeconfig = cluster.kubeconfig
export const oidcProviderArn = cluster.oidcProviderArn
export const eventsTableName = eventsTable.name
export const eventsTableArn = eventsTable.arn
export const eventsTableStreamArn = eventsTable.streamArn
export const opensearchEndpoint = opensearch.endpoint
export const opensearchArn = opensearch.arn
export const documentDbEndpoint = docdb.cluster.endpoint
export const documentDbConnectionString = docdb.connectionString
export const documentDbPort = docdb.port
export const ecrRepositoryUrls = Object.fromEntries(Object.entries(ecrRepositories).map(([k, v]) => [k, v.repositoryUrl]))
export const eventNotificationsQueueUrl = queues.eventNotificationsQueue.url
export const eventNotificationsQueueArn = queues.eventNotificationsQueue.arn
export const eventNotificationsDlqArn = queues.eventNotificationsDlq.arn
export const fanoutLambdaArn = fanoutLambda.lambda.arn
export const frontendBucketName = frontendBucket.bucket
export const cloudfrontDistributionId = cloudfront.distribution.id
export const cloudfrontDomainName = cloudfront.distribution.domainName
export const route53ZoneId = dns.zone?.zoneId
export const route53RecordFqdn = dns.record?.fqdn
export const prometheusRoleArn = prometheus.role.arn
export const grafanaRoleArn = grafana.role.arn
export const eksNodeRoleArn = eksNodeRole.arn
