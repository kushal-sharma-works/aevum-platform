import { onMounted, onUnmounted, ref, watch, type Ref } from 'vue'

export function useInfiniteScroll(
	onLoadMore: () => Promise<void> | void,
	isLoading: Ref<boolean>,
	hasMore: Ref<boolean>
) {
	const sentinel = ref<HTMLElement | null>(null)
	let observer: IntersectionObserver | null = null

	const observe = () => {
		if (!sentinel.value) {
			return
		}

		observer = new IntersectionObserver(async (entries) => {
			const first = entries[0]
			if (!first?.isIntersecting || isLoading.value || !hasMore.value) {
				return
			}
			await onLoadMore()
		})

		observer.observe(sentinel.value)
	}

	onMounted(observe)

	onUnmounted(() => {
		observer?.disconnect()
	})

	watch(sentinel, () => {
		observer?.disconnect()
		observe()
	})

	return { sentinel }
}
