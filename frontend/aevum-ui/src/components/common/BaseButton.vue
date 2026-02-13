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

const classes = computed(() => {
	const base = 'rounded px-3 py-2 text-sm font-medium transition disabled:opacity-60'
	const map = {
		primary: 'bg-blue-600 text-white hover:bg-blue-500',
		secondary: 'bg-slate-700 text-white hover:bg-slate-600',
		danger: 'bg-rose-600 text-white hover:bg-rose-500',
		ghost: 'bg-transparent text-slate-200 hover:bg-slate-800'
	}

	return `${base} ${map[props.variant]}`
})
</script>

<template>
	<button :type="type" :disabled="disabled" :class="classes">
		<slot />
	</button>
</template>
