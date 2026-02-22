import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/client', () => ({
	default: {
		get: vi.fn(),
		post: vi.fn()
	}
}))

import client from '@/api/client'
import { eventTimelineApi } from '@/api/eventTimeline'

describe('eventTimelineApi', () => {
	it('maps methods to expected urls and params', async () => {
		vi.mocked(client.post).mockResolvedValueOnce({ data: { eventId: 'e1' } })
		vi.mocked(client.post).mockResolvedValueOnce({ data: { accepted: 1, rejected: 0 } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { data: [], nextCursor: null, hasMore: false } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { eventId: 'e1' } })
		vi.mocked(client.post).mockResolvedValueOnce({})

		await eventTimelineApi.ingestEvent({ streamId: 's1', eventType: 'Created' } as any)
		await eventTimelineApi.batchIngest({ events: [] } as any)
		await eventTimelineApi.getStreamEvents('s1', 'c1', 10)
		await eventTimelineApi.getEvent('e1')
		await eventTimelineApi.triggerReplay({ streamId: 's1', from: 'a', to: 'b', speed: 1 } as any)

		expect(client.post).toHaveBeenNthCalledWith(1, '/api/events/events', { streamId: 's1', eventType: 'Created' })
		expect(client.post).toHaveBeenNthCalledWith(2, '/api/events/events/batch', { events: [] })
		expect(client.get).toHaveBeenNthCalledWith(1, '/api/events/streams/s1/events', {
			params: { cursor: 'c1', limit: 10 }
		})
		expect(client.get).toHaveBeenNthCalledWith(2, '/api/events/events/e1')
		expect(client.post).toHaveBeenNthCalledWith(3, '/api/events/admin/replay', {
			streamId: 's1',
			from: 'a',
			to: 'b',
			speed: 1
		})
	})
})
