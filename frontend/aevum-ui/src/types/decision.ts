export interface Decision {
	readonly decisionId: string
	readonly eventId: string
	readonly streamId: string
	readonly ruleId: string
	readonly ruleVersion: number
	readonly input: Record<string, unknown>
	readonly output: Record<string, unknown>
	readonly trace: DecisionTrace
	readonly status: DecisionStatus
	readonly deterministicHash: string
	readonly evaluatedAt: string
	readonly eventOccurredAt: string
}

export type DecisionStatus = 'Evaluated' | 'Skipped' | 'Error'

export interface DecisionTrace {
	readonly steps: ReadonlyArray<TraceStep>
	readonly durationMs: number
}

export interface TraceStep {
	readonly conditionField: string
	readonly operator: string
	readonly expectedValue: unknown
	readonly actualValue: unknown
	readonly matched: boolean
	readonly reasoning: string
}

export interface EvaluateRequest {
	readonly streamId: string
	readonly eventId: string
	readonly input: Record<string, unknown>
}

export interface EvaluationResponse {
	readonly decision: Decision
}

export interface DecisionQueryParams {
	readonly streamId?: string
	readonly from?: string
	readonly to?: string
	readonly page?: number
	readonly pageSize?: number
}
