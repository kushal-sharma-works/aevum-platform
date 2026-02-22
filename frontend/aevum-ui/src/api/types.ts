import type { AuditTrail } from '@/types/audit'
import type { Decision, DecisionQueryParams, EvaluateRequest, EvaluationResponse } from '@/types/decision'
import type { DiffResult } from '@/types/diff'
import type { BatchIngestRequest, BatchIngestResponse, Event, IngestEventRequest, ReplayRequest, ReplayResponse } from '@/types/event'
import type { Rule, CreateRuleRequest } from '@/types/rule'
import type { TimelineEntry } from '@/types/timeline'
import type { CursorPage, PaginatedResponse } from '@/types/common'

export interface SearchFilters {
	readonly streamId?: string
	readonly type?: string
	readonly from?: string
	readonly to?: string
}

export interface SearchResult {
	readonly kind: 'event' | 'decision' | 'rule'
	readonly id: string
	readonly streamId?: string
	readonly timestamp: string
	readonly snippet: string
}

export interface TimelineParams {
	readonly stream_id?: string
	readonly from?: string
	readonly to?: string
	readonly type?: string
	readonly page?: number
	readonly size?: number
}

export interface CorrelateParams {
	readonly rule_id?: string
	readonly event_id?: string
	readonly stream_id?: string
}

export interface CorrelationResult {
	readonly items: ReadonlyArray<{ readonly id: string; readonly relation: string }>
}

export type DiffParams =
	| {
			readonly stream_id: string
			readonly t1: string
			readonly t2: string
		}
	| {
			readonly rule_id: string
			readonly v1: number
			readonly v2: number
		}

export type EventTimelineApi = {
	ingestEvent: (req: IngestEventRequest) => Promise<Event>
	batchIngest: (req: BatchIngestRequest) => Promise<BatchIngestResponse>
	getStreamEvents: (streamId: string, cursor?: string, limit?: number) => Promise<CursorPage<Event>>
	getEvent: (eventId: string) => Promise<Event>
	triggerReplay: (req: ReplayRequest) => Promise<ReplayResponse>
}

export type DecisionEngineApi = {
	evaluate: (req: EvaluateRequest) => Promise<EvaluationResponse>
	getDecisions: (params: DecisionQueryParams) => Promise<PaginatedResponse<Decision>>
	getDecision: (id: string) => Promise<Decision>
	getDecisionsByRule: (ruleId: string) => Promise<ReadonlyArray<Decision>>
	createRule: (req: CreateRuleRequest) => Promise<Rule>
	getRules: () => Promise<ReadonlyArray<Rule>>
	getRuleVersions: (ruleId: string) => Promise<ReadonlyArray<Rule>>
	deactivateRuleVersion: (ruleId: string, version: number) => Promise<void>
}

export type QueryAuditApi = {
	search: (query: string, filters?: SearchFilters) => Promise<PaginatedResponse<SearchResult>>
	getTimeline: (params: TimelineParams) => Promise<PaginatedResponse<TimelineEntry>>
	correlate: (params: CorrelateParams) => Promise<CorrelationResult>
	diff: (params: DiffParams) => Promise<DiffResult>
	getAuditTrail: (decisionId: string) => Promise<AuditTrail>
}
