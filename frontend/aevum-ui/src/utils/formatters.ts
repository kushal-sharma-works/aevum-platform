import { format } from 'date-fns'

export function formatTimestamp(iso: string): string {
	const date = new Date(iso)
	if (Number.isNaN(date.getTime())) {
		return 'Invalid date'
	}
	return format(date, 'yyyy-MM-dd HH:mm:ss XXX')
}

export function formatDuration(ms: number): string {
	if (ms < 1000) {
		return `${ms}ms`
	}
	if (ms < 60_000) {
		return `${(ms / 1000).toFixed(1)}s`
	}
	const minutes = Math.floor(ms / 60_000)
	const seconds = Math.round((ms % 60_000) / 1000)
	return `${minutes}m ${seconds}s`
}

export function formatNumber(n: number): string {
	return new Intl.NumberFormat().format(n)
}

export function truncate(s: string, max: number): string {
	if (s.length <= max) {
		return s
	}
	return `${s.slice(0, Math.max(0, max - 1))}…`
}
