<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
	defineProps<{
		variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
		disabled?: boolean
		type?: 'button' | 'submit' | 'reset'
	}>(),
	{
		variant: 'primary',
		disabled: false,
		type: 'button'
	}
)

const color = computed(() => {
	const map = {
		primary: 'primary',
		secondary: 'secondary',
		danger: 'negative',
		ghost: 'primary'
	}

	return map[props.variant]
})
</script>

<template>
	<button
		:type="type"
		:disabled="disabled"
		:class="[
			'p-button p-component',
			variant === 'secondary' ? 'p-button-secondary' : '',
			variant === 'danger' ? 'p-button-danger' : '',
			variant === 'ghost' ? 'p-button-text' : ''
		]"
		style="text-transform: none"
	>
		<slot />
	</button>
</template>
