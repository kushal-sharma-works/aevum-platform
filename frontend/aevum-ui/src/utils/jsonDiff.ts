import type { DiffEntry } from '@/types/diff'

const isObject = (value: unknown): value is Record<string, unknown> => {
	return typeof value === 'object' && value !== null && !Array.isArray(value)
}

export function diffJson(a: unknown, b: unknown, path = ''): DiffEntry[] {
	if (Object.is(a, b)) {
		return []
	}

	if (Array.isArray(a) && Array.isArray(b)) {
		const max = Math.max(a.length, b.length)
		const diffs: DiffEntry[] = []
		for (let i = 0; i < max; i += 1) {
			const nextPath = `${path}[${i}]`
			if (i >= a.length) {
				diffs.push({ path: nextPath, type: 'added', newValue: b[i] })
			} else if (i >= b.length) {
				diffs.push({ path: nextPath, type: 'removed', oldValue: a[i] })
			} else {
				diffs.push(...diffJson(a[i], b[i], nextPath))
			}
		}
		return diffs
	}

	if (isObject(a) && isObject(b)) {
		const keys = new Set([...Object.keys(a), ...Object.keys(b)])
		const diffs: DiffEntry[] = []
		for (const key of keys) {
			const nextPath = path ? `${path}.${key}` : key
			if (!(key in a)) {
				diffs.push({ path: nextPath, type: 'added', newValue: b[key] })
			} else if (!(key in b)) {
				diffs.push({ path: nextPath, type: 'removed', oldValue: a[key] })
			} else {
				diffs.push(...diffJson(a[key], b[key], nextPath))
			}
		}
		return diffs
	}

	return [
		{
			path: path || '$',
			type: 'changed',
			oldValue: a,
			newValue: b
		}
	]
}
