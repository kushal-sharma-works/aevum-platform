export interface Event {
	readonly eventId: string
	readonly streamId: string
	readonly sequenceNumber: number
	readonly eventType: string
	readonly payload: Record<string, unknown>
	readonly metadata: Record<string, string>
	readonly idempotencyKey: string
	readonly occurredAt: string
	readonly ingestedAt: string
	readonly schemaVersion: number
}

export interface IngestEventRequest {
	readonly streamId: string
	readonly eventType: string
	readonly payload: Record<string, unknown>
	readonly metadata?: Record<string, string>
	readonly idempotencyKey?: string
	readonly occurredAt?: string
}

export interface BatchIngestRequest {
	readonly events: ReadonlyArray<IngestEventRequest>
}

export interface BatchIngestResponse {
	readonly ingested: number
	readonly failed: number
}

export interface ReplayRequest {
	readonly streamId: string
	readonly from: string
	readonly to: string
	readonly speed?: number
}
