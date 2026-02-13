import { describe, expect, it } from 'vitest'
import { formatDuration, formatNumber, truncate } from '@/utils/formatters'

describe('formatters', () => {
	it('formats duration', () => {
		expect(formatDuration(42)).toBe('42ms')
	})

	it('formats number', () => {
		expect(formatNumber(1000)).toContain('1')
	})

	it('truncates', () => {
		expect(truncate('abcdef', 3)).toBe('ab…')
	})
})
