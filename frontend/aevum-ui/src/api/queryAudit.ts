import client from './client'
import { normalizeAuditTrail, normalizePaginatedResponse, normalizeTimelineResponse } from './normalizers'
import type {
	CorrelateParams,
	CorrelationResult,
	DiffParams,
	QueryAuditApi,
	SearchFilters,
	SearchResult,
	TimelineParams
} from './types'
import type { AuditTrail } from '@/types/audit'
import type { PaginatedResponse } from '@/types/common'
import type { DiffResult } from '@/types/diff'
import type { TimelineEntry } from '@/types/timeline'

export const queryAuditApi: QueryAuditApi = {
	async search(query: string, filters?: SearchFilters): Promise<PaginatedResponse<SearchResult>> {
		const mappedFilters = filters
			? {
				...filters,
				stream_id: filters.streamId,
				type: filters.type === 'decision' ? 'decisions' : filters.type === 'event' ? 'events' : filters.type
			}
			: undefined

		const { data } = await client.get<PaginatedResponse<SearchResult>>('/api/query/search', {
			params: {
				q: query,
				...mappedFilters
			}
		})
		return normalizePaginatedResponse(data, (item) => item as SearchResult)
	},
	async getTimeline(params: TimelineParams): Promise<PaginatedResponse<TimelineEntry>> {
		const { data } = await client.get<PaginatedResponse<TimelineEntry>>('/api/query/timeline', {
			params
		})
		return normalizeTimelineResponse(data)
	},
	async correlate(params: CorrelateParams): Promise<CorrelationResult> {
		const { data } = await client.get<CorrelationResult>('/api/query/correlate', {
			params
		})
		return data
	},
	async diff(params: DiffParams): Promise<DiffResult> {
		const { data } = await client.get<DiffResult>('/api/query/diff', {
			params
		})
		return data
	},
	async getAuditTrail(decisionId: string): Promise<AuditTrail> {
		const { data } = await client.get<AuditTrail>(`/api/query/audit/${decisionId}`)
		return normalizeAuditTrail(data)
	}
}
