import { ref } from 'vue'
import { defineStore } from 'pinia'
import { decisionEngineApi } from '@/api/decisionEngine'
import type { CreateRuleRequest, Rule } from '@/types/rule'

export const useRulesStore = defineStore('rules', () => {
	const rules = ref<Rule[]>([])
	const selectedRule = ref<Rule | null>(null)
	const ruleVersions = ref<Record<string, ReadonlyArray<Rule>>>({})
	const isLoading = ref(false)
	const error = ref<string | null>(null)

	async function fetchRules(): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			rules.value = [...(await decisionEngineApi.getRules())]
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch rules'
		} finally {
			isLoading.value = false
		}
	}

	async function fetchRule(ruleId: string): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			const allRules = await decisionEngineApi.getRules()
			selectedRule.value = allRules.find((rule) => rule.ruleId === ruleId) ?? null
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch rule'
		} finally {
			isLoading.value = false
		}
	}

	async function fetchRuleVersions(ruleId: string): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			const versions = await decisionEngineApi.getRuleVersions(ruleId)
			ruleVersions.value = {
				...ruleVersions.value,
				[ruleId]: versions
			}
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch rule versions'
		} finally {
			isLoading.value = false
		}
	}

	async function createRule(payload: CreateRuleRequest): Promise<Rule | null> {
		isLoading.value = true
		error.value = null
		try {
			const created = await decisionEngineApi.createRule(payload)
			rules.value = [created, ...rules.value]
			return created
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to create rule'
			return null
		} finally {
			isLoading.value = false
		}
	}

	async function deactivateRuleVersion(ruleId: string, version: number): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			await decisionEngineApi.deactivateRuleVersion(ruleId, version)
			await fetchRuleVersions(ruleId)
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to deactivate rule version'
		} finally {
			isLoading.value = false
		}
	}

	return {
		rules,
		selectedRule,
		ruleVersions,
		isLoading,
		error,
		fetchRules,
		fetchRule,
		fetchRuleVersions,
		createRule,
		deactivateRuleVersion
	}
})
