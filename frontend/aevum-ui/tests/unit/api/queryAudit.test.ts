import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/client', () => ({
	default: {
		get: vi.fn()
	}
}))

import client from '@/api/client'
import { queryAuditApi } from '@/api/queryAudit'

describe('queryAuditApi', () => {
	it('forwards params to all query-audit endpoints', async () => {
		vi.mocked(client.get).mockResolvedValueOnce({ data: { data: [], meta: { page: 1, pageSize: 20, totalCount: 0, totalPages: 1 } } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { data: [], meta: { page: 1, pageSize: 20, totalCount: 0, totalPages: 1 } } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { items: [] } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { before: null, after: null, changes: [] } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { decisionId: 'd1', steps: [] } })

		await queryAuditApi.search('foo', { type: 'event', streamId: 's1' })
		await queryAuditApi.getTimeline({ stream_id: 's1', page: 1, size: 50 })
		await queryAuditApi.correlate({ stream_id: 's1' })
		await queryAuditApi.diff({ stream_id: 's1', t1: '1', t2: '2' })
		await queryAuditApi.getAuditTrail('d1')

		expect(client.get).toHaveBeenNthCalledWith(1, '/api/query/search', {
			params: { query: 'foo', type: 'event', streamId: 's1' }
		})
		expect(client.get).toHaveBeenNthCalledWith(2, '/api/query/timeline', {
			params: { stream_id: 's1', page: 1, size: 50 }
		})
		expect(client.get).toHaveBeenNthCalledWith(3, '/api/query/correlate', { params: { stream_id: 's1' } })
		expect(client.get).toHaveBeenNthCalledWith(4, '/api/query/diff', {
			params: { stream_id: 's1', t1: '1', t2: '2' }
		})
		expect(client.get).toHaveBeenNthCalledWith(5, '/api/query/audit/d1')
	})
})
