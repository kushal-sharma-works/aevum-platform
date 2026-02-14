import { ref } from 'vue'
import { describe, expect, it, vi } from 'vitest'
import { useDebounce } from '@/composables/useDebounce'

describe('useDebounce', () => {
	it('delays updates', () => {
		vi.useFakeTimers()
		const input = ref('a')
		const output = useDebounce(input, 100)
		input.value = 'b'
		expect(output.value).toBe('a')
		vi.advanceTimersByTime(100)
		expect(output.value).toBe('b')
		vi.useRealTimers()
	})
})
