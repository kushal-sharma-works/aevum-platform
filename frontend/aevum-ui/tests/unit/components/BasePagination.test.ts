import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import BasePagination from '@/components/common/BasePagination.vue'

describe('BasePagination', () => {
	it('emits next', async () => {
		const wrapper = mount(BasePagination, { props: { hasMore: true } })
		await wrapper.findAll('button')[1]?.trigger('click')
		expect(wrapper.emitted('next')).toBeTruthy()
	})
})
