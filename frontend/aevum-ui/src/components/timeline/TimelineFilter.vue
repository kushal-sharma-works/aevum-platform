<script setup lang="ts">
const streamId = defineModel<string>('streamId', { required: true })
const type = defineModel<string>('type', { required: true })
</script>

<template>
	<div class="timeline-filter">
		<label class="filter-field">
			<span>Stream ID</span>
			<input v-model="streamId" class="p-inputtext p-component" />
		</label>
		<label class="filter-field">
			<span>Type</span>
			<input v-model="type" class="p-inputtext p-component" />
		</label>
	</div>
</template>

<style scoped>
.timeline-filter {
	display: grid;
	grid-template-columns: repeat(2, minmax(0, 1fr));
	gap: 0.75rem;
}

@media (max-width: 900px) {
	.timeline-filter {
		grid-template-columns: 1fr;
	}
}

.filter-field {
	display: grid;
	gap: 0.35rem;
}

.filter-field span {
	font-size: 0.85rem;
	color: var(--p-text-muted-color);
}
</style>
