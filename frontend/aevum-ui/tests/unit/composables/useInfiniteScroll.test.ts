import { defineComponent, h, nextTick, ref } from 'vue'
import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'

type ObserverCb = (entries: Array<{ isIntersecting: boolean }>) => void

let latestCallback: ObserverCb | null = null

class MockIntersectionObserver {
	constructor(cb: ObserverCb) {
		latestCallback = cb
	}

	observe(): void {}
	disconnect(): void {}
}

describe('useInfiniteScroll', () => {
	beforeEach(() => {
		latestCallback = null
		vi.stubGlobal('IntersectionObserver', MockIntersectionObserver)
	})

	it('loads when sentinel intersects and loading gates allow', async () => {
		const onLoadMore = vi.fn().mockResolvedValue(undefined)
		const isLoading = ref(false)
		const hasMore = ref(true)

		const Host = defineComponent({
			setup() {
				const { sentinel } = useInfiniteScroll(onLoadMore, isLoading, hasMore)
				return () => h('div', [h('div', { ref: sentinel })])
			}
		})

		const wrapper = mount(Host)
		await nextTick()
		expect(latestCallback).toBeTypeOf('function')

		await latestCallback?.([{ isIntersecting: true }])
		expect(onLoadMore).toHaveBeenCalledTimes(1)

		isLoading.value = true
		await latestCallback?.([{ isIntersecting: true }])
		expect(onLoadMore).toHaveBeenCalledTimes(1)

		hasMore.value = false
		isLoading.value = false
		await latestCallback?.([{ isIntersecting: true }])
		expect(onLoadMore).toHaveBeenCalledTimes(1)

		wrapper.unmount()
	})
})
