import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it, vi } from 'vitest'
import { decisionEngineApi } from '@/api/decisionEngine'
import { useDecisionsStore } from '@/stores/decisions'

describe('decisions store', () => {
	it('fetches decisions', async () => {
		setActivePinia(createPinia())
		vi.spyOn(decisionEngineApi, 'getDecisions').mockResolvedValue({
			data: [],
			meta: { page: 1, pageSize: 20, totalCount: 0, totalPages: 1 }
		})
		const store = useDecisionsStore()
		await store.fetchDecisions({ page: 1 })
		expect(store.decisions.length).toBe(0)
	})
})
