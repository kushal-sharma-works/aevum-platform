import { computed, ref, type Ref } from 'vue'

interface VirtualItem<T> {
	readonly index: number
	readonly item: T
}

export function useVirtualList<T>(items: Ref<T[]>, itemHeight: number, containerHeight: number) {
	const scrollTop = ref(0)
	const buffer = 5

	const startIndex = computed(() => Math.max(0, Math.floor(scrollTop.value / itemHeight) - buffer))
	const visibleCount = computed(() => Math.ceil(containerHeight / itemHeight) + buffer * 2)
	const endIndex = computed(() => Math.min(items.value.length, startIndex.value + visibleCount.value))

	const visibleItems = computed<ReadonlyArray<VirtualItem<T>>>(() => {
		return items.value.slice(startIndex.value, endIndex.value).map((item, offset) => ({
			index: startIndex.value + offset,
			item
		}))
	})

	const containerStyle = computed(() => ({
		height: `${containerHeight}px`,
		overflow: 'auto'
	}))

	const listStyle = computed(() => ({
		height: `${items.value.length * itemHeight}px`,
		position: 'relative' as const
	}))

	function scrollTo(index: number): void {
		scrollTop.value = Math.max(0, index * itemHeight)
	}

	return {
		visibleItems,
		containerStyle,
		listStyle,
		scrollTop,
		startIndex,
		itemHeight,
		scrollTo
	}
}
