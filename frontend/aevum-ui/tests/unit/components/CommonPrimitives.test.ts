import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import BaseBadge from '@/components/common/BaseBadge.vue'
import BaseCard from '@/components/common/BaseCard.vue'
import BaseCodeBlock from '@/components/common/BaseCodeBlock.vue'
import BaseEmptyState from '@/components/common/BaseEmptyState.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import BaseSpinner from '@/components/common/BaseSpinner.vue'
import BaseToast from '@/components/common/BaseToast.vue'

describe('common primitives', () => {
	it('renders static primitives', () => {
		expect(mount(BaseBadge, { props: { label: 'ok', tone: 'success' } }).text()).toContain('ok')
		expect(mount(BaseCodeBlock, { props: { code: 'const x = 1' } }).text()).toContain('const x = 1')
		expect(mount(BaseEmptyState, { props: { title: 'No data', description: 'Try later' } }).text()).toContain('No data')
		expect(mount(BaseToast, { props: { message: 'Saved', type: 'success' } }).text()).toContain('Saved')
		expect(mount(BaseSpinner).find('[aria-label="loading"]').exists()).toBe(true)
	})

	it('supports slots and v-model components', async () => {
		const card = mount(BaseCard, {
			slots: {
				default: '<p>inside card</p>'
			}
		})
		expect(card.text()).toContain('inside card')

		const input = mount(BaseInput, {
			props: {
				modelValue: 'abc',
				label: 'Name',
				'onUpdate:modelValue': () => undefined
			}
		})
		await input.find('input').setValue('next')
		expect(input.emitted('update:modelValue')?.[0]).toEqual(['next'])

		const select = mount(BaseSelect, {
			props: {
				modelValue: 'a',
				options: [
					{ label: 'A', value: 'a' },
					{ label: 'B', value: 'b' }
				],
				'onUpdate:modelValue': () => undefined
			}
		})
		await select.find('select').setValue('b')
		expect(select.emitted('update:modelValue')?.[0]).toEqual(['b'])
	})
})
