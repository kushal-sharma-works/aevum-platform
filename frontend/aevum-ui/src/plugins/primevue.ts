import type { App, Component, PropType } from 'vue'
import { defineComponent, h } from 'vue'
import { RouterLink, type RouteLocationRaw, useRouter } from 'vue-router'
import PrimeVue from 'primevue/config'
import Aura from '@primeuix/themes/aura'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Badge from 'primevue/badge'
import Dialog from 'primevue/dialog'
import Card from 'primevue/card'
import Divider from 'primevue/divider'
import Message from 'primevue/message'
import ProgressSpinner from 'primevue/progressspinner'
import ProgressBar from 'primevue/progressbar'
import Slider from 'primevue/slider'

interface QColumn {
	name: string
	label: string
	field: string
}

function toPrimeSeverity(color?: string): string | undefined {
	switch (color) {
		case 'primary':
			return 'primary'
		case 'secondary':
			return 'secondary'
		case 'negative':
		case 'danger':
			return 'danger'
		case 'positive':
			return 'success'
		case 'warning':
			return 'warn'
		case 'info':
			return 'info'
		default:
			return undefined
	}
}

function toPrimeIcon(icon?: string): string | undefined {
	if (!icon) {
		return undefined
	}

	const mapping: Record<string, string> = {
		close: 'pi-times',
		remove: 'pi-minus',
		add: 'pi-plus'
	}

	return `pi ${mapping[icon] ?? 'pi-circle-fill'}`
}

const QBtn = defineComponent({
	name: 'QBtn',
	props: {
		label: String,
		color: String,
		disable: Boolean,
		disabled: Boolean,
		type: { type: String, default: 'button' },
		flat: Boolean,
		outline: Boolean,
		noCaps: Boolean,
		icon: String,
		loading: Boolean,
		to: [String, Object] as PropType<RouteLocationRaw>
	},
	emits: ['click'],
	setup(props, { slots, emit, attrs }) {
		const router = useRouter()
		const onClick = async (event: Event) => {
			if (props.to) {
				await router.push(props.to)
			}
			emit('click', event)
		}

		return () =>
			h(Button as Component, {
				...attrs,
				label: props.label,
				severity: toPrimeSeverity(props.color),
				disabled: props.disable || props.disabled,
				type: props.type,
				text: props.flat,
				outlined: props.outline,
				loading: props.loading,
				icon: toPrimeIcon(props.icon),
				style: props.noCaps ? 'text-transform: none;' : undefined,
				onClick
			}, slots)
	}
})

const QInput = defineComponent({
	name: 'QInput',
	props: {
		modelValue: [String, Number],
		placeholder: String,
		type: { type: String, default: 'text' },
		autocomplete: String,
		label: String
	},
	emits: ['update:modelValue'],
	setup(props, { emit, attrs }) {
		return () =>
			h('label', { class: 'q-field-compat' }, [
				props.label ? h('span', { class: 'q-field-label-compat' }, props.label) : null,
				h(InputText as Component, {
					...attrs,
					modelValue: props.modelValue,
					placeholder: props.placeholder,
					type: props.type,
					autocomplete: props.autocomplete,
					'onUpdate:modelValue': (value: string) => emit('update:modelValue', value)
				})
			])
	}
})

const QSelect = defineComponent({
	name: 'QSelect',
	props: {
		modelValue: [String, Number],
		options: { type: Array, default: () => [] },
		label: String,
		emitValue: Boolean,
		mapOptions: Boolean
	},
	emits: ['update:modelValue'],
	setup(props, { emit, attrs }) {
		return () =>
			h('label', { class: 'q-field-compat' }, [
				props.label ? h('span', { class: 'q-field-label-compat' }, props.label) : null,
				h(Select as Component, {
					...attrs,
					modelValue: props.modelValue,
					options: props.options,
					optionLabel: 'label',
					optionValue: 'value',
					'onUpdate:modelValue': (value: string) => emit('update:modelValue', value)
				})
			])
	}
})

const QBadge = defineComponent({
	name: 'QBadge',
	props: {
		color: String
	},
	setup(props, { slots }) {
		return () => h(Badge as Component, { severity: toPrimeSeverity(props.color), value: String(slots.default?.()[0]?.children ?? '') })
	}
})

const QBanner = defineComponent({
	name: 'QBanner',
	setup(_, { slots, attrs }) {
		const className = String(attrs.class ?? '')
		const severity = className.includes('negative')
			? 'error'
			: className.includes('warning') || className.includes('amber') || className.includes('orange')
				? 'warn'
				: className.includes('positive') || className.includes('teal')
					? 'success'
					: 'info'

		return () => h(Message as Component, { severity }, slots)
	}
})

const QCard = defineComponent({
	name: 'QCard',
	setup(_, { slots, attrs }) {
		return () => h(Card as Component, { ...attrs }, { content: () => slots.default?.() })
	}
})

const QCardSection = defineComponent({
	name: 'QCardSection',
	setup(_, { slots, attrs }) {
		return () => h('div', { ...attrs, class: ['p-4', attrs.class] }, slots.default?.())
	}
})

const QDialog = defineComponent({
	name: 'QDialog',
	props: {
		modelValue: Boolean
	},
	emits: ['hide'],
	setup(props, { slots, emit }) {
		return () =>
			h(Dialog as Component, {
				visible: props.modelValue,
				modal: true,
				dismissableMask: true,
				'onUpdate:visible': (visible: boolean) => {
					if (!visible) emit('hide')
				}
			}, {
				default: () => slots.default?.()
			})
	}
})

const QSeparator = defineComponent({
	name: 'QSeparator',
	setup() {
		return () => h(Divider as Component)
	}
})

const QSpinner = defineComponent({
	name: 'QSpinner',
	setup() {
		return () => h(ProgressSpinner as Component, { style: 'width: 20px; height: 20px' })
	}
})

const QLinearProgress = defineComponent({
	name: 'QLinearProgress',
	props: {
		value: { type: Number, default: 0 }
	},
	setup(props) {
		return () => h(ProgressBar as Component, { value: Math.min(100, Math.max(0, props.value * 100)) })
	}
})

const QSlider = defineComponent({
	name: 'QSlider',
	props: {
		modelValue: Number,
		min: Number,
		max: Number
	},
	emits: ['update:modelValue'],
	setup(props, { emit }) {
		return () =>
			h(Slider as Component, {
				modelValue: props.modelValue,
				min: props.min,
				max: props.max,
				'onUpdate:modelValue': (value: number) => emit('update:modelValue', value)
			})
	}
})

const QTable = defineComponent({
	name: 'QTable',
	props: {
		rows: { type: Array, default: () => [] },
		columns: { type: Array, default: () => [] },
		rowKey: { type: String, default: 'id' }
	},
	setup(props, { slots }) {
		return () =>
			h('table', { class: 'p-datatable-table q-table-compat' }, [
				h('thead', [
					h('tr', (props.columns as QColumn[]).map((column) => h('th', { key: column.name }, column.label)))
				]),
				h(
					'tbody',
					(props.rows as Record<string, unknown>[]).map((row, rowIndex) =>
						h('tr', { key: String(row[props.rowKey] ?? rowIndex) }, (props.columns as QColumn[]).map((column) => {
							const slotName = `body-cell-${column.name}`
							const scoped = slots[slotName]?.({ row, col: column, value: row[column.field], props: { row, col: column } })
							return scoped?.length ? scoped : h('td', { key: column.name }, String(row[column.field] ?? ''))
						}))
					)
				)
			])
	}
})

const QTd = defineComponent({
	name: 'QTd',
	setup(_, { slots }) {
		return () => h('td', slots.default?.())
	}
})

const QItem = defineComponent({
	name: 'QItem',
	props: {
		to: [String, Object] as PropType<RouteLocationRaw>
	},
	setup(props, { slots }) {
		return () =>
			props.to
				? h(RouterLink as Component, { to: props.to, class: 'q-item-compat p-button p-button-text' }, () => slots.default?.())
				: h('div', { class: 'q-item-compat' }, slots.default?.())
	}
})

const QItemSection = defineComponent({
	name: 'QItemSection',
	setup(_, { slots }) {
		return () => h('span', { class: 'q-item-section-compat' }, slots.default?.())
	}
})

const QItemLabel = defineComponent({
	name: 'QItemLabel',
	setup(_, { slots }) {
		return () => h('span', slots.default?.())
	}
})

const QIcon = defineComponent({
	name: 'QIcon',
	props: {
		name: { type: String, default: 'circle' }
	},
	setup(props) {
		return () => h('i', { class: `pi pi-circle-fill`, style: 'font-size: 0.6rem' })
	}
})

const QShellBlock = defineComponent({
	name: 'QShellBlock',
	setup(_, { slots }) {
		return () => h('div', slots.default?.())
	}
})

const QToolbar = defineComponent({
	name: 'QToolbar',
	setup(_, { slots }) {
		return () => h('div', { class: 'q-toolbar-compat' }, slots.default?.())
	}
})

const QToolbarTitle = defineComponent({
	name: 'QToolbarTitle',
	setup(_, { slots }) {
		return () => h('div', { class: 'q-toolbar-title-compat' }, slots.default?.())
	}
})

const QList = defineComponent({
	name: 'QList',
	setup(_, { slots }) {
		return () => h('div', { class: 'q-list-compat' }, slots.default?.())
	}
})

const aliases: Record<string, Component> = {
	'q-btn': QBtn,
	'q-input': QInput,
	'q-select': QSelect,
	'q-badge': QBadge,
	'q-banner': QBanner,
	'q-card': QCard,
	'q-card-section': QCardSection,
	'q-dialog': QDialog,
	'q-separator': QSeparator,
	'q-spinner': QSpinner,
	'q-linear-progress': QLinearProgress,
	'q-slider': QSlider,
	'q-table': QTable,
	'q-td': QTd,
	'q-item': QItem,
	'q-item-section': QItemSection,
	'q-item-label': QItemLabel,
	'q-icon': QIcon,
	'q-form': defineComponent({ name: 'QForm', emits: ['submit'], setup(_, { slots, emit }) { return () => h('form', { onSubmit: (event: Event) => emit('submit', event) }, slots.default?.()) } }),
	'q-layout': QShellBlock,
	'q-header': defineComponent({ setup(_, { slots }) { return () => h('header', slots.default?.()) } }),
	'q-footer': defineComponent({ setup(_, { slots }) { return () => h('footer', slots.default?.()) } }),
	'q-drawer': defineComponent({ setup(_, { slots }) { return () => h('aside', { class: 'q-drawer-compat' }, slots.default?.()) } }),
	'q-page-container': QShellBlock,
	'q-page': defineComponent({ setup(_, { slots }) { return () => h('section', slots.default?.()) } }),
	'q-toolbar': QToolbar,
	'q-toolbar-title': QToolbarTitle,
	'q-list': QList
}

export function installPrimeVue(app: App): void {
	app.use(PrimeVue, {
		theme: {
			preset: Aura
		}
	})

	Object.entries(aliases).forEach(([name, component]) => {
		app.component(name, component)
	})
}
