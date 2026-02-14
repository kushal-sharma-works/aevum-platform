export type ConditionOperator =
	| 'Eq'
	| 'NotEq'
	| 'Gt'
	| 'Gte'
	| 'Lt'
	| 'Lte'
	| 'Contains'
	| 'NotContains'
	| 'Regex'
	| 'In'
	| 'NotIn'

export type ActionType = 'Approve' | 'Reject' | 'Flag' | 'Escalate' | 'Transform'

export interface RuleCondition {
	readonly field: string
	readonly operator: ConditionOperator
	readonly value: unknown
}

export interface RuleAction {
	readonly actionType: ActionType
	readonly parameters: Record<string, unknown>
}

export interface Rule {
	readonly ruleId: string
	readonly name: string
	readonly description: string
	readonly version: number
	readonly conditions: ReadonlyArray<RuleCondition>
	readonly actions: ReadonlyArray<RuleAction>
	readonly isActive: boolean
	readonly createdAt: string
	readonly createdBy: string
}

export interface CreateRuleRequest {
	readonly name: string
	readonly description: string
	readonly conditions: ReadonlyArray<RuleCondition>
	readonly actions: ReadonlyArray<RuleAction>
}
