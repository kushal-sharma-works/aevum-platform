import { ref, watchEffect, type Ref } from 'vue'

export function useDebounce<T>(value: Ref<T>, delay: number): Ref<T> {
	const debounced = ref(value.value) as Ref<T>
	let timer: number | undefined

	watchEffect((onCleanup) => {
		timer = window.setTimeout(() => {
			debounced.value = value.value
		}, delay)

		onCleanup(() => {
			if (timer !== undefined) {
				window.clearTimeout(timer)
			}
		})
	})

	return debounced
}
