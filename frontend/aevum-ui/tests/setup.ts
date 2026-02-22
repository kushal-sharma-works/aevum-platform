import { defineComponent } from 'vue'
import { config } from '@vue/test-utils'

config.global.stubs = {
	'q-btn': defineComponent({
		props: {
			disable: Boolean,
			label: String,
			type: { type: String, default: 'button' }
		},
		emits: ['click'],
		template: `<button :type="type" :disabled="disable" @click="$emit('click', $event)"><slot />{{ label }}</button>`
	}),
	'q-input': defineComponent({
		props: {
			modelValue: { type: [String, Number], default: '' },
			placeholder: String
		},
		emits: ['update:modelValue'],
		template:
			'<input :value="modelValue" :placeholder="placeholder" @input="$emit(\'update:modelValue\', $event.target.value)" />'
	}),
	'q-select': defineComponent({
		props: {
			modelValue: { type: [String, Number], default: '' },
			options: { type: Array, default: () => [] }
		},
		emits: ['update:modelValue'],
		template: `<select :value="modelValue" @change="$emit('update:modelValue', $event.target.value)">
			<option v-for="option in options" :key="option.value ?? option" :value="option.value ?? option">{{ option.label ?? option }}</option>
		</select>`
	}),
	'q-dialog': defineComponent({
		props: {
			modelValue: Boolean
		},
		template: '<div v-if="modelValue"><slot /></div>'
	}),
	'q-card': defineComponent({ template: '<div><slot /></div>' }),
	'q-card-section': defineComponent({ template: '<div><slot /></div>' }),
	'q-banner': defineComponent({ template: '<div><slot /></div>' }),
	'q-spinner': defineComponent({ template: '<div class="spinner-stub" />' }),
	'q-separator': defineComponent({ template: '<hr />' }),
	'q-layout': defineComponent({ template: '<div><slot /></div>' }),
	'q-header': defineComponent({ template: '<header><slot /></header>' }),
	'q-toolbar': defineComponent({ template: '<div><slot /></div>' }),
	'q-toolbar-title': defineComponent({ template: '<div><slot /></div>' }),
	'q-badge': defineComponent({ template: '<span><slot /></span>' }),
	'q-drawer': defineComponent({ template: '<aside><slot /></aside>' }),
	'q-list': defineComponent({ template: '<div><slot /></div>' }),
	'q-item': defineComponent({ template: '<div><slot /></div>' }),
	'q-item-label': defineComponent({ template: '<div><slot /></div>' }),
	'q-item-section': defineComponent({ template: '<div><slot /></div>' }),
	'q-icon': defineComponent({ template: '<i />' }),
	'q-page-container': defineComponent({ template: '<main><slot /></main>' }),
	'q-page': defineComponent({ template: '<section><slot /></section>' }),
	'q-footer': defineComponent({ template: '<footer><slot /></footer>' }),
	'q-linear-progress': defineComponent({ template: '<div />' })
}

if (!window.matchMedia) {
	Object.defineProperty(window, 'matchMedia', {
		writable: true,
		value: (query: string) => ({
			matches: false,
			media: query,
			onchange: null,
			addEventListener: () => undefined,
			removeEventListener: () => undefined,
			dispatchEvent: () => false,
			addListener: () => undefined,
			removeListener: () => undefined
		})
	})
}

if (!('ResizeObserver' in window)) {
	class ResizeObserverMock {
		observe(): void {}
		unobserve(): void {}
		disconnect(): void {}
	}

	Object.defineProperty(window, 'ResizeObserver', {
		writable: true,
		value: ResizeObserverMock
	})
}

if (!navigator.clipboard) {
	Object.defineProperty(navigator, 'clipboard', {
		value: {
			writeText: async () => undefined
		}
	})
}