<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>

<template>
	<section class="empty-state">
		<div class="empty-title">{{ title }}</div>
		<div v-if="description" class="empty-description">{{ description }}</div>
	</section>
</template>

<style scoped>
.empty-state {
	border: 1px solid var(--p-content-border-color);
	background: var(--p-content-background);
	border-radius: 0.75rem;
	padding: 1.25rem;
	text-align: center;
}

.empty-title {
	font-weight: 600;
	font-size: 1.05rem;
}

.empty-description {
	margin-top: 0.25rem;
	color: var(--p-text-muted-color);
	font-size: 0.9rem;
}
</style>
