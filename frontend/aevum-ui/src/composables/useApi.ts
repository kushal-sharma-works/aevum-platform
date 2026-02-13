import { ref, type Ref } from 'vue'

export function useApi<T>(fetcher: () => Promise<T>) {
	const data = ref<T | null>(null) as Ref<T | null>
	const error = ref<Error | null>(null)
	const isLoading = ref(false)

	async function execute(): Promise<T | null> {
		isLoading.value = true
		error.value = null
		try {
			const result = await fetcher()
			data.value = result
			return result
		} catch (err) {
			error.value = err instanceof Error ? err : new Error('Unknown error')
			return null
		} finally {
			isLoading.value = false
		}
	}

	return { data, error, isLoading, execute }
}
