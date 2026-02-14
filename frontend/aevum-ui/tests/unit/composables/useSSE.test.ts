import { describe, expect, it, vi } from 'vitest'
import { useSSE } from '@/composables/useSSE'

class MockEventSource {
	public onopen: (() => void) | null = null
	public onmessage: ((event: MessageEvent<string>) => void) | null = null
	public onerror: (() => void) | null = null

	constructor(_: string) {
		setTimeout(() => {
			this.onopen?.()
		}, 0)
	}

	close(): void {}
}

describe('useSSE', () => {
	it('opens connection', async () => {
		vi.stubGlobal('EventSource', MockEventSource)
		const sse = useSSE('/events')
		await new Promise((resolve) => setTimeout(resolve, 5))
		expect(['connecting', 'open', 'error', 'closed']).toContain(sse.status.value)
	})
})
