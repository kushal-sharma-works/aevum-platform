import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'
import TimelineEntry from '@/components/timeline/TimelineEntry.vue'

describe('TimelineEntry', () => {
	it('renders entry kind', () => {
		const wrapper = mount(TimelineEntry, {
			props: {
				entry: {
					kind: 'event',
					timestamp: new Date().toISOString(),
					item: {
						eventId: 'e1',
						streamId: 's1',
						sequenceNumber: 1,
						eventType: 'Created',
						payload: {},
						metadata: {},
						idempotencyKey: 'k',
						occurredAt: new Date().toISOString(),
						ingestedAt: new Date().toISOString(),
						schemaVersion: 1
					}
				}
			}
		})
		expect(wrapper.text()).toContain('event')
	})
})
