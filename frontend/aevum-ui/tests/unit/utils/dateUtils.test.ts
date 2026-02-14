import { describe, expect, it } from 'vitest'
import { getPresetRange, isWithinRange } from '@/utils/dateUtils'

describe('dateUtils', () => {
	it('returns preset range', () => {
		const range = getPresetRange('Last 1h')
		expect(range.from.getTime()).toBeLessThan(range.to.getTime())
	})

	it('checks range inclusion', () => {
		const from = new Date(Date.now() - 1000)
		const to = new Date(Date.now() + 1000)
		expect(isWithinRange(new Date(), from, to)).toBe(true)
	})
})
