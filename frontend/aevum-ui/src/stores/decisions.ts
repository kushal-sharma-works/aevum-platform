import { ref } from 'vue'
import { defineStore } from 'pinia'
import { decisionEngineApi } from '@/api/decisionEngine'
import type { Decision, DecisionQueryParams, EvaluateRequest, EvaluationResponse } from '@/types/decision'
import type { PaginatedResponse } from '@/types/common'

export const useDecisionsStore = defineStore('decisions', () => {
	const decisions = ref<Decision[]>([])
	const selectedDecision = ref<Decision | null>(null)
	const isLoading = ref(false)
	const error = ref<string | null>(null)
	const pageMeta = ref<PaginatedResponse<Decision>['meta'] | null>(null)

	async function fetchDecisions(filters: DecisionQueryParams): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			const rules = await decisionEngineApi.getRules()
			const decisionGroups = await Promise.all(rules.map((rule) => decisionEngineApi.getDecisionsByRule(rule.ruleId).catch(() => [])))
			const allDecisions = decisionGroups.flat()

			const filtered = allDecisions.filter((decision) => {
				if (filters.streamId && decision.streamId !== filters.streamId) {
					return false
				}

				const evaluatedAt = Date.parse(decision.evaluatedAt)
				if (filters.from && Number.isFinite(evaluatedAt) && evaluatedAt < Date.parse(filters.from)) {
					return false
				}
				if (filters.to && Number.isFinite(evaluatedAt) && evaluatedAt > Date.parse(filters.to)) {
					return false
				}

				return true
			})

			const sorted = filtered.sort((left, right) => Date.parse(right.evaluatedAt) - Date.parse(left.evaluatedAt))
			const pageSize = filters.pageSize ?? 20
			const page = filters.page ?? 1
			const start = (page - 1) * pageSize
			const end = start + pageSize
			decisions.value = sorted.slice(start, end)
			pageMeta.value = {
				page,
				pageSize,
				totalCount: sorted.length,
				totalPages: Math.max(1, Math.ceil(sorted.length / pageSize))
			}
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch decisions'
		} finally {
			isLoading.value = false
		}
	}

	async function fetchDecision(decisionId: string): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			selectedDecision.value = await decisionEngineApi.getDecision(decisionId)
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch decision'
		} finally {
			isLoading.value = false
		}
	}

	async function evaluate(request: EvaluateRequest): Promise<EvaluationResponse | null> {
		isLoading.value = true
		error.value = null
		try {
			const response = await decisionEngineApi.evaluate(request)
			return response
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to evaluate decision'
			return null
		} finally {
			isLoading.value = false
		}
	}

	return {
		decisions,
		selectedDecision,
		pageMeta,
		isLoading,
		error,
		fetchDecisions,
		fetchDecision,
		evaluate
	}
})
