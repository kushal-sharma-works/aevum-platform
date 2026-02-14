import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { eventTimelineApi } from '@/api/eventTimeline'
import { useSSE } from '@/composables/useSSE'
import type { Decision } from '@/types/decision'
import type { Event } from '@/types/event'

type ReplayItem = Event | Decision

export const useReplayStore = defineStore('replay', () => {
	const isReplaying = ref(false)
	const replayProgress = ref(0)
	const replayedEvents = ref<ReplayItem[]>([])
	const replaySpeed = ref<1 | 2 | 5 | 10>(1)
	const streamId = ref('')
	const from = ref('')
	const to = ref('')

	const sseUrl = computed(() => `/api/query/replay/stream?stream_id=${encodeURIComponent(streamId.value)}`)
	const sse = useSSE(sseUrl)

	watch(
		() => sse.data.value,
		(payload) => {
			if (!isReplaying.value || payload === null || typeof payload !== 'object') {
				return
			}
			replayedEvents.value = [payload as ReplayItem, ...replayedEvents.value]
			replayProgress.value = Math.min(100, replayProgress.value + 2)
		}
	)

	async function startReplay(nextStreamId: string, nextFrom: string, nextTo: string): Promise<void> {
		streamId.value = nextStreamId
		from.value = nextFrom
		to.value = nextTo
		replayProgress.value = 0
		replayedEvents.value = []
		isReplaying.value = true

		await eventTimelineApi.triggerReplay({
			streamId: nextStreamId,
			from: nextFrom,
			to: nextTo,
			speed: replaySpeed.value
		})
	}

	function stopReplay(): void {
		isReplaying.value = false
		replayProgress.value = 0
		sse.close()
	}

	return {
		isReplaying,
		replayProgress,
		replayedEvents,
		replaySpeed,
		startReplay,
		stopReplay
	}
})
