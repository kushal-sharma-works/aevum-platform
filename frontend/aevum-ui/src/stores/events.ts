import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { eventTimelineApi } from '@/api/eventTimeline'
import type { IngestEventRequest, Event } from '@/types/event'

export const useEventsStore = defineStore('events', () => {
	const events = ref<Event[]>([])
	const selectedEvent = ref<Event | null>(null)
	const isLoading = ref(false)
	const error = ref<string | null>(null)
	const nextCursor = ref<string | null>(null)
	const hasMore = computed(() => nextCursor.value !== null)

	async function fetchStreamEvents(streamId: string, cursor?: string, limit = 100): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			const page = await eventTimelineApi.getStreamEvents(streamId, cursor, limit)
			events.value = cursor ? [...events.value, ...page.data] : [...page.data]
			nextCursor.value = page.nextCursor
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch events'
		} finally {
			isLoading.value = false
		}
	}

	async function fetchEvent(eventId: string): Promise<void> {
		isLoading.value = true
		error.value = null
		selectedEvent.value = null
		try {
			selectedEvent.value = await eventTimelineApi.getEvent(eventId)
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to fetch event'
		} finally {
			isLoading.value = false
		}
	}

	async function ingestEvent(payload: IngestEventRequest): Promise<void> {
		isLoading.value = true
		error.value = null
		try {
			const created = await eventTimelineApi.ingestEvent(payload)
			events.value = [created, ...events.value]
		} catch (err) {
			error.value = err instanceof Error ? err.message : 'Failed to ingest event'
		} finally {
			isLoading.value = false
		}
	}

	return {
		events,
		selectedEvent,
		isLoading,
		error,
		nextCursor,
		hasMore,
		fetchStreamEvents,
		fetchEvent,
		ingestEvent
	}
})
