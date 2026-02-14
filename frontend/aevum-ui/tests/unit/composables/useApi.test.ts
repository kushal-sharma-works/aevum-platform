import { describe, expect, it } from 'vitest'
import { useApi } from '@/composables/useApi'

describe('useApi', () => {
	it('returns data on execute', async () => {
		const api = useApi(async () => 42)
		const result = await api.execute()
		expect(result).toBe(42)
		expect(api.data.value).toBe(42)
	})
})
