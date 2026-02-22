<script setup lang="ts">
defineProps<{ step: { readonly timestamp: string; readonly message: string } }>()
</script>

<template>
	<section class="audit-step">
		<div class="audit-time">{{ step.timestamp }}</div>
		<div>{{ step.message }}</div>
	</section>
</template>

<style scoped>
.audit-step {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.6rem 0.75rem;
}

.audit-time {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
	margin-bottom: 0.1rem;
}
</style>
