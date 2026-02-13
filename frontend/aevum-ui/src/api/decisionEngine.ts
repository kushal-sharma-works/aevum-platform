import client from './client'
import type { DecisionEngineApi } from './types'
import type { PaginatedResponse } from '@/types/common'
import type { Decision, DecisionQueryParams, EvaluateRequest, EvaluationResponse } from '@/types/decision'
import type { CreateRuleRequest, Rule } from '@/types/rule'

export const decisionEngineApi: DecisionEngineApi = {
	async evaluate(req: EvaluateRequest): Promise<EvaluationResponse> {
		const { data } = await client.post<EvaluationResponse>('/api/decisions/evaluate', req)
		return data
	},
	async getDecisions(params: DecisionQueryParams): Promise<PaginatedResponse<Decision>> {
		const { data } = await client.get<PaginatedResponse<Decision>>('/api/decisions/decisions', { params })
		return data
	},
	async getDecision(id: string): Promise<Decision> {
		const { data } = await client.get<Decision>(`/api/decisions/decisions/${id}`)
		return data
	},
	async createRule(req: CreateRuleRequest): Promise<Rule> {
		const { data } = await client.post<Rule>('/api/decisions/rules', req)
		return data
	},
	async getRules(): Promise<ReadonlyArray<Rule>> {
		const { data } = await client.get<ReadonlyArray<Rule>>('/api/decisions/rules')
		return data
	},
	async getRuleVersions(ruleId: string): Promise<ReadonlyArray<Rule>> {
		const { data } = await client.get<ReadonlyArray<Rule>>(`/api/decisions/rules/${ruleId}/versions`)
		return data
	},
	async deactivateRuleVersion(ruleId: string, version: number): Promise<void> {
		await client.put(`/api/decisions/rules/${ruleId}/versions/${version}/deactivate`)
	}
}
