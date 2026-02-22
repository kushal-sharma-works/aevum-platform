<script setup lang="ts">
import type { Rule } from '@/types/rule'

defineProps<{ rule: Rule }>()
</script>

<template>
	<section class="rule-detail">
		<div class="rule-name">{{ rule.name }} v{{ rule.version }}</div>
		<div class="rule-description">{{ rule.description }}</div>
	</section>
</template>

<style scoped>
.rule-detail {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.rule-name {
	font-size: 1rem;
	font-weight: 600;
}

.rule-description {
	margin-top: 0.2rem;
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
