import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import BaseButton from '@/components/common/BaseButton.vue'

describe('BaseButton', () => {
	it('renders slot content', () => {
		const wrapper = mount(BaseButton, { slots: { default: 'Click' } })
		expect(wrapper.text()).toContain('Click')
	})
})
