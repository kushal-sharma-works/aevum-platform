import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import BaseModal from '@/components/common/BaseModal.vue'

describe('BaseModal', () => {
	it('renders when open', () => {
		mount(BaseModal, { props: { open: true, title: 'T' }, attachTo: document.body })
		expect(document.body.textContent).toContain('T')
	})
})
