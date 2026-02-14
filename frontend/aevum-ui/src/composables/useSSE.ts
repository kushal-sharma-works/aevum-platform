import { isRef, onUnmounted, ref, type Ref } from 'vue'

export type SSEStatus = 'connecting' | 'open' | 'closed' | 'error'

export function useSSE(url: string | Ref<string>) {
	const data = ref<unknown>(null)
	const error = ref<string | null>(null)
	const status = ref<SSEStatus>('connecting')
	const source = ref<EventSource | null>(null)

	let retries = 0
	let retryTimer: number | undefined

	const getUrl = () => (isRef(url) ? url.value : url)

	const connect = () => {
		status.value = 'connecting'
		source.value = new EventSource(getUrl())

		source.value.onopen = () => {
			retries = 0
			status.value = 'open'
		}

		source.value.onmessage = (event: MessageEvent<string>) => {
			try {
				data.value = JSON.parse(event.data) as unknown
			} catch {
				data.value = event.data
			}
		}

		source.value.onerror = () => {
			status.value = 'error'
			error.value = 'SSE connection failed'
			source.value?.close()
			const delay = Math.min(30000, 1000 * 2 ** retries)
			retries += 1
			retryTimer = window.setTimeout(connect, delay)
		}
	}

	const close = () => {
		if (retryTimer !== undefined) {
			window.clearTimeout(retryTimer)
		}
		source.value?.close()
		status.value = 'closed'
	}

	connect()
	onUnmounted(close)

	return { data, error, status, close }
}
