import { describe, expect, it } from 'vitest'
import { diffJson } from '@/utils/jsonDiff'

describe('jsonDiff', () => {
	it('detects changed values', () => {
		const diffs = diffJson({ a: 1 }, { a: 2 })
		expect(diffs.length).toBeGreaterThan(0)
	})

	it('detects added keys', () => {
		const diffs = diffJson({}, { a: 1 })
		expect(diffs.some((entry) => entry.type === 'added')).toBe(true)
	})
})
