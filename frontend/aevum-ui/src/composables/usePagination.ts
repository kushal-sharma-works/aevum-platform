import { computed, ref } from 'vue'

export type PaginationMode = 'cursor' | 'page'

export function usePagination(mode: PaginationMode = 'page') {
	const currentPage = ref(1)
	const nextCursor = ref<string | null>(null)
	const cursorStack = ref<string[]>([])

	const hasMore = computed(() => (mode === 'page' ? true : nextCursor.value !== null))

	function goNext(): void {
		if (mode === 'page') {
			currentPage.value += 1
			return
		}

		if (nextCursor.value !== null) {
			cursorStack.value.push(nextCursor.value)
		}
	}

	function goPrev(): void {
		if (mode === 'page') {
			currentPage.value = Math.max(1, currentPage.value - 1)
			return
		}

		cursorStack.value.pop()
		nextCursor.value = cursorStack.value[cursorStack.value.length - 1] ?? null
	}

	function reset(): void {
		currentPage.value = 1
		nextCursor.value = null
		cursorStack.value = []
	}

	return {
		currentPage,
		nextCursor,
		hasMore,
		goNext,
		goPrev,
		reset,
		cursorStack
	}
}
