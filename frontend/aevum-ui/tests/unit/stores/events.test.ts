import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it, vi } from 'vitest'
import { eventTimelineApi } from '@/api/eventTimeline'
import { useEventsStore } from '@/stores/events'

describe('events store', () => {
	it('fetches stream events', async () => {
		setActivePinia(createPinia())
		vi.spyOn(eventTimelineApi, 'getStreamEvents').mockResolvedValue({ data: [], nextCursor: null, hasMore: false })
		const store = useEventsStore()
		await store.fetchStreamEvents('s1')
		expect(store.events.length).toBe(0)
	})
})
