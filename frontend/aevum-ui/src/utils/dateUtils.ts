import { parseISO as parseDateISO, subDays, subHours } from 'date-fns'
import type { TimePreset } from '@/composables/useTimeRange'

export function parseISO(s: string): Date {
	return parseDateISO(s)
}

export function isWithinRange(date: Date, from: Date, to: Date): boolean {
	const time = date.getTime()
	return time >= from.getTime() && time <= to.getTime()
}

export function getPresetRange(preset: TimePreset): { from: Date; to: Date } {
	const to = new Date()
	switch (preset) {
		case 'Last 1h':
			return { from: subHours(to, 1), to }
		case 'Last 24h':
			return { from: subHours(to, 24), to }
		case 'Last 7d':
			return { from: subDays(to, 7), to }
		case 'Last 30d':
			return { from: subDays(to, 30), to }
		case 'Custom':
			return { from: subHours(to, 1), to }
	}
}
