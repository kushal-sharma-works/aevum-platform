<script setup lang="ts">
const props = withDefaults(
	defineProps<{
		page?: number
		totalPages?: number
		hasMore?: boolean
	}>(),
	{
		page: 1,
		totalPages: 1,
		hasMore: false
	}
)

const emit = defineEmits<{
	(e: 'next'): void
	(e: 'prev'): void
}>()
</script>

<template>
	<div class="flex items-center gap-2">
		<button class="rounded bg-slate-800 px-3 py-1" :disabled="page <= 1" @click="emit('prev')">Prev</button>
		<span class="text-sm text-slate-400">Page {{ page }} / {{ totalPages }}</span>
		<button class="rounded bg-slate-800 px-3 py-1" :disabled="!hasMore && page >= totalPages" @click="emit('next')">
			Next
		</button>
	</div>
</template>
