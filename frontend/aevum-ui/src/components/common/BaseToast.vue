<script setup lang="ts">
defineProps<{ message: string; type?: 'success' | 'error' | 'warning' | 'info' }>()
</script>

<template>
	<div
		:class="[
			'rounded px-3 py-2 text-sm shadow',
			type === 'success'
				? 'bg-emerald-700'
				: type === 'error'
					? 'bg-rose-700'
					: type === 'warning'
						? 'bg-amber-700'
						: 'bg-blue-700'
		]"
	>
		{{ message }}
	</div>
</template>
