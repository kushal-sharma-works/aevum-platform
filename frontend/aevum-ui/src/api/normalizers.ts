import type { AuditTrail } from '@/types/audit'
import type { CursorPage, PaginatedResponse } from '@/types/common'
import type { Decision, DecisionTrace, TraceStep } from '@/types/decision'
import type { Event } from '@/types/event'
import type { Rule } from '@/types/rule'
import type { TimelineEntry } from '@/types/timeline'

type AnyRecord = Record<string, unknown>

function asRecord(value: unknown): AnyRecord {
	return value !== null && typeof value === 'object' ? (value as AnyRecord) : {}
}

function asArray(value: unknown): unknown[] {
	return Array.isArray(value) ? value : []
}

function asString(value: unknown, fallback = ''): string {
	return typeof value === 'string' ? value : fallback
}

function asNumber(value: unknown, fallback = 0): number {
	return typeof value === 'number' && Number.isFinite(value) ? value : fallback
}

function asBoolean(value: unknown, fallback = false): boolean {
	return typeof value === 'boolean' ? value : fallback
}

function normalizeDecisionStatus(value: unknown): Decision['status'] {
	if (typeof value === 'string') {
		if (value === 'Evaluated' || value === 'Skipped' || value === 'Error') {
			return value
		}
	}

	if (typeof value === 'number') {
		switch (value) {
			case 1:
				return 'Evaluated'
			case 2:
				return 'Skipped'
			default:
				return 'Error'
		}
	}

	return 'Evaluated'
}

export function normalizeEvent(input: unknown): Event {
	const record = asRecord(input)
	return {
		eventId: asString(record.eventId ?? record.event_id),
		streamId: asString(record.streamId ?? record.stream_id),
		sequenceNumber: asNumber(record.sequenceNumber ?? record.sequence_number),
		eventType: asString(record.eventType ?? record.event_type),
		payload: asRecord(record.payload),
		metadata: asRecord(record.metadata) as Record<string, string>,
		idempotencyKey: asString(record.idempotencyKey ?? record.idempotency_key),
		occurredAt: asString(record.occurredAt ?? record.occurred_at),
		ingestedAt: asString(record.ingestedAt ?? record.ingested_at),
		schemaVersion: asNumber(record.schemaVersion ?? record.schema_version, 1)
	}
}

function normalizeTraceStep(input: unknown, index: number): TraceStep {
	const record = asRecord(input)
	return {
		conditionField: asString(record.conditionField ?? record.field, `step_${index}`),
		operator: asString(record.operator),
		expectedValue: record.expectedValue ?? record.expected ?? null,
		actualValue: record.actualValue ?? record.actual ?? null,
		matched: asBoolean(record.matched),
		reasoning: asString(record.reasoning ?? record.message)
	}
}

function normalizeDecisionTrace(input: unknown): DecisionTrace {
	const record = asRecord(input)
	const steps = asArray(record.steps).map((step, index) => normalizeTraceStep(step, index))
	return {
		steps,
		durationMs: asNumber(record.durationMs ?? record.duration_ms)
	}
}

export function normalizeDecision(input: unknown): Decision {
	const record = asRecord(input)
	const matchedConditions = asArray(record.matchedConditions ?? record.matched_conditions)
	const traceSteps = matchedConditions.map((condition, index) => {
		const item = asRecord(condition)
		return normalizeTraceStep(
			{
				conditionField: item.field,
				operator: item.operator,
				expectedValue: item.expectedValue ?? item.value,
				actualValue: item.actualValue,
				matched: item.matched,
				reasoning: item.reasoning
			},
			index
		)
	})

	const trace =
		traceSteps.length > 0 || record.trace
			? normalizeDecisionTrace({
				steps: traceSteps.length > 0 ? traceSteps : asRecord(record.trace).steps,
				durationMs: record.evaluationDurationMs ?? asRecord(record.trace).durationMs,
				duration_ms: asRecord(record.trace).duration_ms
			})
			: normalizeDecisionTrace({ durationMs: record.evaluationDurationMs })

	return {
		decisionId: asString(record.decisionId ?? record.decision_id ?? record.id),
		eventId: asString(record.eventId ?? record.event_id ?? record.requestId ?? record.request_id),
		streamId: asString(record.streamId ?? record.stream_id),
		ruleId: asString(record.ruleId ?? record.rule_id),
		ruleVersion: asNumber(record.ruleVersion ?? record.rule_version ?? record.version),
		input: asRecord(record.input ?? record.inputContext ?? record.input_context),
		output: asRecord(record.output ?? record.outputData ?? record.output_data),
		trace,
		status: normalizeDecisionStatus(record.status),
		deterministicHash: asString(record.deterministicHash ?? record.deterministic_hash),
		evaluatedAt: asString(record.evaluatedAt ?? record.evaluated_at),
		eventOccurredAt: asString(record.eventOccurredAt ?? record.event_occurred_at)
	}
}

export function normalizeRule(input: unknown): Rule {
	const record = asRecord(input)
	const ruleId = asString(record.ruleId ?? record.rule_id ?? record.id)
	return {
		ruleId,
		name: asString(record.name),
		description: asString(record.description),
		version: asNumber(record.version, 1),
		conditions: asArray(record.conditions) as Rule['conditions'],
		actions: asArray(record.actions).map((action, index) => {
			const actionRecord = asRecord(action)
			return {
				actionType: asString(actionRecord.actionType ?? actionRecord.type, 'Approve') as Rule['actions'][number]['actionType'],
				parameters: asRecord(actionRecord.parameters)
			}
		}) as Rule['actions'],
		isActive: asBoolean(record.isActive, true),
		createdAt: asString(record.createdAt ?? record.created_at),
		createdBy: asString(record.createdBy ?? record.created_by, 'system')
	}
}

function normalizeTimelineEntry(input: unknown): TimelineEntry {
	const record = asRecord(input)
	const inferredType = record.event_id ?? record.eventId ? 'event' : record.decision_id ?? record.decisionId ?? record.ruleId ? 'decision' : ''
	const type = asString(record.kind ?? record.type ?? inferredType)
	const timestamp = asString(record.timestamp ?? record.occurred_at ?? record.ingested_at ?? record.evaluatedAt ?? record.evaluated_at)
	if (type === 'event') {
		return {
			kind: 'event',
			timestamp,
			item: normalizeEvent(record.item ?? record.event ?? record)
		}
	}

	return {
		kind: 'decision',
		timestamp,
		item: normalizeDecision(record.item ?? record.decision ?? record)
	}
}

export function normalizePaginatedResponse<T>(input: unknown, mapper: (value: unknown) => T): PaginatedResponse<T> {
	const record = asRecord(input)
	const dataSource = asArray(record.data ?? record.items ?? record.hits)
	const data = dataSource.map((item) => mapper(item))

	const page = asNumber(record.page ?? asRecord(record.meta).page, 1)
	const pageSize = asNumber(record.pageSize ?? asRecord(record.meta).pageSize, data.length)
	const totalCount = asNumber(record.totalCount ?? record.total ?? asRecord(record.meta).totalCount, data.length)
	const totalPages = asNumber(record.totalPages ?? asRecord(record.meta).totalPages, Math.max(1, Math.ceil(totalCount / Math.max(1, pageSize))))

	return {
		data,
		meta: {
			page,
			pageSize,
			totalCount,
			totalPages
		}
	}
}

export function normalizeCursorPage<T>(input: unknown, mapper: (value: unknown) => T): CursorPage<T> {
	const record = asRecord(input)
	const dataSource = asArray(record.data ?? record.items ?? record.events)
	const nextCursor = record.nextCursor ?? record.next_cursor ?? null

	return {
		data: dataSource.map((item) => mapper(item)),
		nextCursor: typeof nextCursor === 'string' || nextCursor === null ? nextCursor : null,
		hasMore: asBoolean(record.hasMore ?? record.has_more, Boolean(nextCursor))
	}
}

export function normalizeTimelineResponse(input: unknown): PaginatedResponse<TimelineEntry> {
	return normalizePaginatedResponse(input, normalizeTimelineEntry)
}

export function normalizeAuditTrail(input: unknown): AuditTrail {
	const record = asRecord(input)
	const steps = asArray(record.timeline ?? record.steps).map((item) => {
		const step = asRecord(item)
		return {
			timestamp: asString(step.timestamp),
			message: asString(step.message)
		}
	})

	const chainRecord = asRecord(record.chain)

	return {
		decisionId: asString(record.decisionId ?? record.decision_id),
		streamId: asString(record.streamId ?? record.stream_id),
		timeline: steps,
		chain: {
			nodes: asArray(chainRecord.nodes).map((node) => {
				const nodeRecord = asRecord(node)
				return {
					id: asString(nodeRecord.id),
					type: asString(nodeRecord.type),
					label: asString(nodeRecord.label)
				}
			}),
			edges: asArray(chainRecord.edges).map((edge) => {
				const edgeRecord = asRecord(edge)
				return {
					from: asString(edgeRecord.from),
					to: asString(edgeRecord.to),
					relation: asString(edgeRecord.relation)
				}
			})
		}
	}
}