import { describe, expect, it } from 'vitest'
import { usePagination } from '@/composables/usePagination'

describe('usePagination', () => {
	it('increments page mode', () => {
		const pager = usePagination('page')
		pager.goNext()
		expect(pager.currentPage.value).toBe(2)
	})
})
