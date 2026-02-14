export interface PaginatedResponse<T> {
	readonly data: ReadonlyArray<T>
	readonly meta: {
		readonly page: number
		readonly pageSize: number
		readonly totalCount: number
		readonly totalPages: number
	}
}

export interface CursorPage<T> {
	readonly data: ReadonlyArray<T>
	readonly nextCursor: string | null
	readonly hasMore: boolean
}

export interface ApiError {
	readonly type: string
	readonly title: string
	readonly status: number
	readonly detail: string
	readonly traceId?: string
}
