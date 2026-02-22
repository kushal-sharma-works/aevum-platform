<script setup lang="ts">
const value = defineModel<number>({ required: true })
</script>

<template>
	<div class="scrubber-wrap">
		<input v-model.number="value" type="range" min="0" max="100" class="scrubber" />
		<span class="scrubber-value">{{ value }}%</span>
	</div>
</template>

<style scoped>
.scrubber-wrap {
	display: grid;
	gap: 0.35rem;
}

.scrubber {
	width: 100%;
	accent-color: var(--p-primary-color);
}

.scrubber-value {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
