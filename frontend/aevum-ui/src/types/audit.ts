export interface CausalChainNode {
	readonly id: string
	readonly type: string
	readonly label: string
}

export interface CausalChainEdge {
	readonly from: string
	readonly to: string
	readonly relation: string
}

export interface CausalChain {
	readonly nodes: ReadonlyArray<CausalChainNode>
	readonly edges: ReadonlyArray<CausalChainEdge>
}

export interface AuditTrail {
	readonly decisionId: string
	readonly streamId: string
	readonly timeline: ReadonlyArray<{
		readonly timestamp: string
		readonly message: string
	}>
	readonly chain: CausalChain
}
