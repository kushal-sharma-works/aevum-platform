import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import ConditionBuilder from '@/components/rules/ConditionBuilder.vue'

describe('ConditionBuilder', () => {
	it('adds condition row', async () => {
		const wrapper = mount(ConditionBuilder, { props: { modelValue: [] } as never })
		await wrapper.find('button').trigger('click')
		expect(wrapper.emitted()).toBeTruthy()
	})
})
