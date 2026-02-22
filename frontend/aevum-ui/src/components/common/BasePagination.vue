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
	<div class="row items-center q-col-gutter-sm">
		<div>
			<button type="button" class="p-button p-button-outlined" :disabled="page <= 1" @click="emit('prev')">Prev</button>
		</div>
		<div class="text-caption text-grey-7">Page {{ page }} / {{ totalPages }}</div>
		<div>
			<button type="button" class="p-button p-button-outlined" :disabled="!hasMore && page >= totalPages" @click="emit('next')">Next</button>
		</div>
	</div>
</template>
