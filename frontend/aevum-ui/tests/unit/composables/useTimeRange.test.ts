import { describe, expect, it } from 'vitest'
import { useTimeRange } from '@/composables/useTimeRange'

describe('useTimeRange', () => {
	it('updates range for preset values', () => {
		const range = useTimeRange()
		range.setPreset('Last 24h')

		expect(range.preset.value).toBe('Last 24h')
		expect(range.isValid.value).toBe(true)
		expect(range.to.value.getTime()).toBeGreaterThan(range.from.value.getTime())
	})

	it('ignores invalid custom ranges and accepts valid ones', () => {
		const range = useTimeRange()
		const originalFrom = range.from.value
		const originalTo = range.to.value

		range.setCustomRange(new Date('2026-01-02T00:00:00Z'), new Date('2026-01-01T00:00:00Z'))
		expect(range.from.value).toEqual(originalFrom)
		expect(range.to.value).toEqual(originalTo)

		range.setCustomRange(new Date('2026-01-01T00:00:00Z'), new Date('2026-01-02T00:00:00Z'))
		expect(range.preset.value).toBe('Custom')
		expect(range.from.value.toISOString()).toBe('2026-01-01T00:00:00.000Z')
		expect(range.to.value.toISOString()).toBe('2026-01-02T00:00:00.000Z')
	})
})
