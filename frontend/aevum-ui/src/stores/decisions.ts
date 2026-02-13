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
			const result = await decisionEngineApi.getDecisions(filters)
			decisions.value = [...result.data]
			pageMeta.value = result.meta
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
