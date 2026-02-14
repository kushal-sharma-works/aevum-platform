import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import DecisionTrace from '@/components/decisions/DecisionTrace.vue'

describe('DecisionTrace', () => {
	it('renders steps', () => {
		const wrapper = mount(DecisionTrace, {
			props: {
				trace: {
					durationMs: 12,
					steps: [
						{
							conditionField: 'age',
							operator: 'Gt',
							expectedValue: 18,
							actualValue: 21,
							matched: true,
							reasoning: 'ok'
						}
					]
				}
			}
		})
		expect(wrapper.text()).toContain('age')
	})
})
