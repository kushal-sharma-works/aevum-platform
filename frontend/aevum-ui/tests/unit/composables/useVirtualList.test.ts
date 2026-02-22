import { ref } from 'vue'
import { describe, expect, it } from 'vitest'
import { useVirtualList } from '@/composables/useVirtualList'

describe('useVirtualList', () => {
	it('computes visible window and styles', () => {
		const items = ref(Array.from({ length: 100 }, (_, index) => ({ id: index })))
		const list = useVirtualList(items, 20, 100)

		expect(list.containerStyle.value).toEqual({ height: '100px', overflow: 'auto' })
		expect(list.listStyle.value.height).toBe('2000px')
		expect(list.visibleItems.value.length).toBeGreaterThan(0)
		expect(list.visibleItems.value[0]?.index).toBe(0)

		list.scrollTo(10)
		expect(list.scrollTop.value).toBe(200)
		expect(list.startIndex.value).toBe(5)
		expect(list.visibleItems.value[0]?.index).toBe(5)
	})
})
