import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it, vi } from 'vitest'
import { decisionEngineApi } from '@/api/decisionEngine'
import { useRulesStore } from '@/stores/rules'

describe('rules store', () => {
	it('fetches rules', async () => {
		setActivePinia(createPinia())
		vi.spyOn(decisionEngineApi, 'getRules').mockResolvedValue([])
		const store = useRulesStore()
		await store.fetchRules()
		expect(store.rules.length).toBe(0)
	})
})
