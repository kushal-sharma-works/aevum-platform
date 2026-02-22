import { describe, expect, it, vi } from 'vitest'

vi.mock('@/api/client', () => ({
	default: {
		get: vi.fn(),
		post: vi.fn(),
		put: vi.fn()
	}
}))

import client from '@/api/client'
import { decisionEngineApi } from '@/api/decisionEngine'

describe('decisionEngineApi', () => {
	it('calls expected endpoints', async () => {
		vi.mocked(client.post).mockResolvedValueOnce({ data: { decisionId: 'd1' } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { data: [], meta: { page: 1, pageSize: 20, totalCount: 0, totalPages: 1 } } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: { decisionId: 'd1' } })
		vi.mocked(client.post).mockResolvedValueOnce({ data: { ruleId: 'r1' } })
		vi.mocked(client.get).mockResolvedValueOnce({ data: [{ ruleId: 'r1' }] })
		vi.mocked(client.get).mockResolvedValueOnce({ data: [{ ruleId: 'r1', version: 2 }] })
		vi.mocked(client.put).mockResolvedValueOnce({})

		await decisionEngineApi.evaluate({ streamId: 's1', input: {} } as any)
		await decisionEngineApi.getDecisions({ page: 1 } as any)
		await decisionEngineApi.getDecision('d1')
		await decisionEngineApi.createRule({ name: 'rule' } as any)
		await decisionEngineApi.getRules()
		await decisionEngineApi.getRuleVersions('r1')
		await decisionEngineApi.deactivateRuleVersion('r1', 2)

		expect(client.post).toHaveBeenNthCalledWith(1, '/api/decisions/evaluate', { streamId: 's1', input: {} })
		expect(client.get).toHaveBeenNthCalledWith(1, '/api/decisions/decisions', { params: { page: 1 } })
		expect(client.get).toHaveBeenNthCalledWith(2, '/api/decisions/decisions/d1')
		expect(client.post).toHaveBeenNthCalledWith(2, '/api/decisions/rules', { name: 'rule' })
		expect(client.get).toHaveBeenNthCalledWith(3, '/api/decisions/rules')
		expect(client.get).toHaveBeenNthCalledWith(4, '/api/decisions/rules/r1/versions')
		expect(client.put).toHaveBeenCalledWith('/api/decisions/rules/r1/versions/2/deactivate')
	})
})
