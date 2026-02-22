import { ref } from 'vue'
import { defineStore } from 'pinia'
import { eventTimelineApi } from '@/api/eventTimeline'
import type { Decision } from '@/types/decision'
import type { Event } from '@/types/event'

type ReplayItem = Event | Decision

export const useReplayStore = defineStore('replay', () => {
	const isReplaying = ref(false)
	const isLoading = ref(false)
	const error = ref<string | null>(null)
	const replayProgress = ref(0)
	const replayedEvents = ref<ReplayItem[]>([])
	const replaySpeed = ref<1 | 2 | 5 | 10>(1)
	const streamId = ref('')
	const from = ref('')
	const to = ref('')

	async function startReplay(nextStreamId: string, nextFrom: string, nextTo: string): Promise<void> {
		const streamIdValue = nextStreamId.trim()
		if (!streamIdValue) {
			error.value = 'Stream ID is required'
			return
		}

		if (!nextFrom || !nextTo) {
			error.value = 'From and To timestamps are required'
			return
		}

		isLoading.value = true
		error.value = null
		streamId.value = nextStreamId
		from.value = nextFrom
		to.value = nextTo
		replayProgress.value = 0
		replayedEvents.value = []

		try {
			isReplaying.value = true
			const result = await eventTimelineApi.triggerReplay({
				streamId: streamIdValue,
				from: nextFrom,
				to: nextTo,
				speed: replaySpeed.value
			})
			replayProgress.value = 100
			error.value = result.eventsReplayed >= 0 ? null : 'Replay completed with unknown result'
			isReplaying.value = false
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to start replay'
			isReplaying.value = false
		} finally {
			isLoading.value = false
		}
	}

	function stopReplay(): void {
		isReplaying.value = false
		replayProgress.value = 0
	}

	return {
		isReplaying,
		isLoading,
		error,
		replayProgress,
		replayedEvents,
		replaySpeed,
		startReplay,
		stopReplay
	}
})
