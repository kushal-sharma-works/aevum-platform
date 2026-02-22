<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>

<template>
	<section class="page-container">
		<header>
			<div class="page-title">{{ title }}</div>
			<div v-if="description" class="page-description">{{ description }}</div>
		</header>
		<hr class="page-separator" />
		<div>
			<slot />
		</div>
	</section>
</template>

<style scoped>
.page-container {
	display: grid;
	gap: 0.85rem;
}

.page-title {
	font-size: 1.25rem;
	font-weight: 600;
}

.page-description {
	margin-top: 0.15rem;
	color: var(--p-text-muted-color);
	font-size: 0.95rem;
}

.page-separator {
	margin: 0;
	border: 0;
	border-top: 1px solid var(--p-content-border-color);
}
</style>
