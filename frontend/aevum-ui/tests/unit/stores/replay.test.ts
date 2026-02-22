import { nextTick, ref } from 'vue'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/api/eventTimeline', () => ({
	eventTimelineApi: {
		triggerReplay: vi.fn()
	}
}))

vi.mock('@/composables/useSSE', () => ({
	useSSE: vi.fn()
}))

import { eventTimelineApi } from '@/api/eventTimeline'
import { useSSE } from '@/composables/useSSE'
import { useReplayStore } from '@/stores/replay'

describe('replay store', () => {
	const sseData = ref<unknown>(null)
	const closeMock = vi.fn()

	beforeEach(() => {
		setActivePinia(createPinia())
		vi.mocked(eventTimelineApi.triggerReplay).mockReset()
		vi.mocked(useSSE).mockReset()
		closeMock.mockReset()
		sseData.value = null

		vi.mocked(useSSE).mockReturnValue({
			data: sseData,
			error: ref<string | null>(null),
			status: ref<'connecting' | 'open' | 'closed' | 'error'>('open'),
			close: closeMock
		})
	})

	it('starts replay and consumes SSE payloads', async () => {
		vi.mocked(eventTimelineApi.triggerReplay).mockResolvedValue(undefined)
		const store = useReplayStore()

		await store.startReplay('stream-1', '2026-01-01', '2026-01-02')

		expect(eventTimelineApi.triggerReplay).toHaveBeenCalledWith({
			streamId: 'stream-1',
			from: '2026-01-01',
			to: '2026-01-02',
			speed: 1
		})
		expect(store.isReplaying).toBe(true)
		expect(store.replayProgress).toBe(0)

		sseData.value = { id: 'evt-1' }
		await nextTick()

		expect(store.replayedEvents).toHaveLength(1)
		expect(store.replayProgress).toBe(2)
	})

	it('stops replay and closes sse', () => {
		const store = useReplayStore()
		store.stopReplay()

		expect(store.isReplaying).toBe(false)
		expect(store.replayProgress).toBe(0)
		expect(closeMock).toHaveBeenCalledTimes(1)
	})
})
