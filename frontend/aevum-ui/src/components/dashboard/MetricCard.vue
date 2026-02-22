<script setup lang="ts">
defineProps<{ label: string; value: string | number }>()
</script>

<template>
	<section class="metric-card">
		<div class="metric-label">{{ label }}</div>
		<div class="metric-value">{{ value }}</div>
	</section>
</template>

<style scoped>
.metric-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.metric-label {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}

.metric-value {
	font-size: 1.2rem;
	font-weight: 600;
}
</style>
