import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { queryAuditApi } from '@/api/queryAudit'
import type { TimelineEntry, TimelineStats } from '@/types/timeline'

export const useTimelineStore = defineStore('timeline', () => {
	const entries = ref<TimelineEntry[]>([])
	const isLoading = ref(false)
	const error = ref<string | null>(null)
	const timeRange = ref({
		from: new Date(Date.now() - 24 * 60 * 60 * 1000),
		to: new Date()
	})

	const stats = computed<TimelineStats>(() => {
		const eventCount = entries.value.filter((entry) => entry.kind === 'event').length
		const decisions = entries.value.filter((entry) => entry.kind === 'decision')
		const decisionCount = decisions.length
		const statusCounts: Record<string, number> = {}

		for (const decision of decisions) {
			const status = decision.item.status
			statusCounts[status] = (statusCounts[status] ?? 0) + 1
		}

		return {
			eventCount,
			decisionCount,
			statusCounts
		}
	})

	async function fetchTimeline(streamId?: string, from?: string, to?: string, type?: string): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			const response = await queryAuditApi.getTimeline({
				stream_id: streamId,
				from,
				to,
				type,
				page: 1,
				size: 200
			})

			entries.value = [...response.data].sort((left, right) => {
				return new Date(left.timestamp).getTime() - new Date(right.timestamp).getTime()
			})
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch timeline'
		} finally {
			isLoading.value = false
		}
	}

	return {
		entries,
		isLoading,
		error,
		timeRange,
		stats,
		fetchTimeline
	}
})
