import client from './client'
import type { EventTimelineApi } from './types'
import type { BatchIngestRequest, BatchIngestResponse, Event, IngestEventRequest, ReplayRequest } from '@/types/event'
import type { CursorPage } from '@/types/common'

export const eventTimelineApi: EventTimelineApi = {
	async ingestEvent(req: IngestEventRequest): Promise<Event> {
		const { data } = await client.post<Event>('/api/events/events', req)
		return data
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
		return data
	},
	async getEvent(eventId: string): Promise<Event> {
		const { data } = await client.get<Event>(`/api/events/events/${eventId}`)
		return data
	},
	async triggerReplay(req: ReplayRequest): Promise<void> {
		await client.post('/api/events/admin/replay', req)
	}
}
