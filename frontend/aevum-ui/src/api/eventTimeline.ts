import client from './client'
import type { EventTimelineApi } from './types'
import { normalizeCursorPage, normalizeEvent } from './normalizers'
import type { BatchIngestRequest, BatchIngestResponse, Event, IngestEventRequest, ReplayRequest, ReplayResponse } from '@/types/event'
import type { CursorPage } from '@/types/common'

export const eventTimelineApi: EventTimelineApi = {
	async ingestEvent(req: IngestEventRequest): Promise<Event> {
		const { data } = await client.post<Event>('/api/events/events', req)
		return normalizeEvent(data)
	},
	async batchIngest(req: BatchIngestRequest): Promise<BatchIngestResponse> {
		const { data } = await client.post<BatchIngestResponse>('/api/events/events/batch', req)
		return data
	},
	async getStreamEvents(streamId: string, cursor?: string, limit?: number): Promise<CursorPage<Event>> {
		const { data } = await client.get<CursorPage<Event>>(`/api/events/streams/${streamId}/events`, {
			params: {
				cursor,
				limit
			}
		})
		return normalizeCursorPage(data, normalizeEvent)
	},
	async getEvent(eventId: string): Promise<Event> {
		const { data } = await client.get<Event | { event?: Event }>(`/api/events/events/${eventId}`)
		const payload = (data as { event?: unknown }).event ?? data
		return normalizeEvent(payload)
	},
	async triggerReplay(req: ReplayRequest): Promise<ReplayResponse> {
		const { data } = await client.post<{ status?: string; events_replayed?: number; eventsReplayed?: number }>('/api/events/admin/replay', {
			stream_id: req.streamId,
			from: req.from,
			to: req.to,
			speed_factor: req.speed
		})

		return {
			status: data.status ?? 'completed',
			eventsReplayed: Number(data.eventsReplayed ?? data.events_replayed ?? 0)
		}
	}
}
