import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it, vi } from 'vitest'
import { queryAuditApi } from '@/api/queryAudit'
import { useTimelineStore } from '@/stores/timeline'

describe('timeline store', () => {
	it('fetches timeline', async () => {
		setActivePinia(createPinia())
		vi.spyOn(queryAuditApi, 'getTimeline').mockResolvedValue({
			data: [],
			meta: { page: 1, pageSize: 20, totalCount: 0, totalPages: 1 }
		})
		const store = useTimelineStore()
		await store.fetchTimeline()
		expect(store.entries.length).toBe(0)
	})
})
