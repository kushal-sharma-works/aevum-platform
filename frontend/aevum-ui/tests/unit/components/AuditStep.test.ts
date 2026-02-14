import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import AuditStep from '@/components/audit/AuditStep.vue'

describe('AuditStep', () => {
	it('renders message', () => {
		const wrapper = mount(AuditStep, {
			props: {
				step: {
					timestamp: new Date().toISOString(),
					message: 'm'
				}
			}
		})
		expect(wrapper.text()).toContain('m')
	})
})
