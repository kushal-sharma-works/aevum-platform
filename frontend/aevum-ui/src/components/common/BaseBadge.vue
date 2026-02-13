<script setup lang="ts">
defineProps<{ label: string; tone?: 'success' | 'warning' | 'error' | 'info' }>()
</script>

<template>
	<span
		:class="[
			'inline-flex rounded px-2 py-0.5 text-xs',
			tone === 'success'
				? 'bg-emerald-900 text-emerald-300'
				: tone === 'warning'
					? 'bg-amber-900 text-amber-300'
					: tone === 'error'
						? 'bg-rose-900 text-rose-300'
						: 'bg-blue-900 text-blue-300'
		]"
	>
		{{ label }}
	</span>
</template>
