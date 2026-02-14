import { InfraConfig } from "../../types"

export interface NodeGroupConfig {
  desiredSize: number
  minSize: number
  maxSize: number
  instanceType: string
  labels: Record<string, string>
}

export function nodeGroupConfig(config: InfraConfig): NodeGroupConfig {
  return {
    desiredSize: config.eksNodeCount,
    minSize: config.eksMinSize,
    maxSize: config.eksMaxSize,
    instanceType: config.eksInstanceType,
    labels: {
      workload: "general",
      environment: config.env
    }
  }
}
