import type { Decision } from './decision'
import type { Event } from './event'

export interface TimeRange {
	readonly from: Date
	readonly to: Date
}

export type TimelineEntry =
	| {
			readonly kind: 'event'
			readonly timestamp: string
			readonly item: Event
		}
	| {
			readonly kind: 'decision'
			readonly timestamp: string
			readonly item: Decision
		}

export interface TimelineStats {
	readonly eventCount: number
	readonly decisionCount: number
	readonly statusCounts: Readonly<Record<string, number>>
}
