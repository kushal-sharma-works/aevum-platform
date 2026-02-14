export interface DiffEntry {
	readonly path: string
	readonly type: 'added' | 'removed' | 'changed'
	readonly oldValue?: unknown
	readonly newValue?: unknown
}

export interface DiffResult {
	readonly entries: ReadonlyArray<DiffEntry>
}
