import client from './client'
import type { DecisionEngineApi } from './types'
import { normalizeDecision, normalizePaginatedResponse, normalizeRule } from './normalizers'
import type { PaginatedResponse } from '@/types/common'
import type { Decision, DecisionQueryParams, EvaluateRequest, EvaluationResponse } from '@/types/decision'
import type { CreateRuleRequest, Rule } from '@/types/rule'

function withoutUndefined<T extends Record<string, unknown>>(value: T): T {
	return Object.fromEntries(Object.entries(value).filter(([, item]) => item !== undefined)) as T
}

function mapOperator(operator: string): number {
	switch (operator) {
		case 'Eq':
			return 0
		case 'NotEq':
			return 1
		case 'Gt':
			return 2
		case 'Gte':
			return 3
		case 'Lt':
			return 4
		case 'Lte':
			return 5
		case 'Contains':
			return 6
		case 'NotContains':
			return 7
		case 'In':
			return 10
		case 'NotIn':
			return 11
		case 'Regex':
			return 12
		default:
			return 0
	}
}

function mapActionType(): number {
	return 0
}

export const decisionEngineApi: DecisionEngineApi = {
	async evaluate(req: EvaluateRequest): Promise<EvaluationResponse> {
		const { data } = await client.post<EvaluationResponse>('/api/decisions/decisions/evaluate', req)
		return {
			decision: normalizeDecision((data as { decision?: unknown }).decision ?? data)
		}
	},
	async getDecisions(params: DecisionQueryParams): Promise<PaginatedResponse<Decision>> {
		const query = withoutUndefined({
			...params,
			stream_id: params.streamId,
			page_size: params.pageSize
		})
		const { data } = await client.get<PaginatedResponse<Decision>>('/api/decisions/decisions', { params: query })
		return normalizePaginatedResponse(data, normalizeDecision)
	},
	async getDecision(id: string): Promise<Decision> {
		const { data } = await client.get<Decision>(`/api/decisions/decisions/${id}`)
		return normalizeDecision(data)
	},
	async createRule(req: CreateRuleRequest): Promise<Rule> {
		const actions = req.actions ?? []
		const payload = actions.length
			? {
				...req,
				priority: 1,
				conditions: req.conditions.map((condition) => ({
					...condition,
					operator: mapOperator(condition.operator)
				})),
				actions: actions.map((action, index) => ({
					type: mapActionType(),
					order: index + 1,
					parameters: {
						...action.parameters,
						actionType: action.actionType
					}
				}))
			}
			: {
				...req,
				priority: 1,
				conditions: req.conditions.map((condition) => ({
					...condition,
					operator: mapOperator(condition.operator)
				}))
			}

		const { data } = await client.post<Rule>('/api/decisions/rules', payload)
		return normalizeRule(data)
	},
	async getRules(): Promise<ReadonlyArray<Rule>> {
		const { data } = await client.get<ReadonlyArray<Rule>>('/api/decisions/rules')
		return data.map((rule) => normalizeRule(rule))
	},
	async getRuleVersions(ruleId: string): Promise<ReadonlyArray<Rule>> {
		const { data } = await client.get<Rule>(`/api/decisions/rules/${ruleId}`)
		return [normalizeRule(data)]
	},
	async getDecisionsByRule(ruleId: string): Promise<ReadonlyArray<Decision>> {
		const { data } = await client.get<ReadonlyArray<Decision>>(`/api/decisions/decisions/rule/${ruleId}`)
		return data.map((decision) => normalizeDecision(decision))
	},
	async deactivateRuleVersion(ruleId: string, version: number): Promise<void> {
		void version
		await client.post(`/api/decisions/rules/${ruleId}/deactivate`)
	}
}
