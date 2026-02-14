import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import { createRouter, createWebHistory } from 'vue-router'
import DecisionCard from '@/components/decisions/DecisionCard.vue'

const router = createRouter({ history: createWebHistory(), routes: [{ path: '/decisions/:decisionId', component: { template: '<div />' } }] })

describe('DecisionCard', () => {
	it('renders decision id', async () => {
		await router.push('/')
		await router.isReady()
		const wrapper = mount(DecisionCard, {
			global: { plugins: [router] },
			props: {
				decision: {
					decisionId: 'd1',
					eventId: 'e1',
					streamId: 's1',
					ruleId: 'r1',
					ruleVersion: 1,
					input: {},
					output: {},
					trace: { steps: [], durationMs: 1 },
					status: 'Evaluated',
					deterministicHash: 'h',
					evaluatedAt: new Date().toISOString(),
					eventOccurredAt: new Date().toISOString()
				}
			}
		})
		expect(wrapper.text()).toContain('d1')
	})
})
