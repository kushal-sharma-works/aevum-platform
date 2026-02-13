import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'

describe('BaseJsonViewer', () => {
	it('shows object keys', () => {
		const wrapper = mount(BaseJsonViewer, { props: { value: { a: 1 } } })
		expect(wrapper.text()).toContain('a')
	})
})
