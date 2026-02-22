<script setup lang="ts">
import type { Rule } from '@/types/rule'

defineProps<{ rules: ReadonlyArray<Rule> }>()
</script>

<template>
	<div class="rule-list">
		<RouterLink
			v-for="rule in rules"
			:key="`${rule.ruleId}-${rule.version}`"
			:to="`/rules/${rule.ruleId}`"
			class="rule-row"
		>
			<div>
				<div class="rule-name">{{ rule.name }}</div>
				<div class="rule-id">{{ rule.ruleId }}</div>
			</div>
			<span class="p-tag p-tag-primary">v{{ rule.version }}</span>
		</RouterLink>
	</div>
</template>

<style scoped>
.rule-list {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	overflow: hidden;
	background: var(--p-content-background);
}

.rule-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 0.75rem;
	padding: 0.75rem;
	text-decoration: none;
	color: inherit;
	border-bottom: 1px solid var(--p-content-border-color);
}

.rule-row:last-child {
	border-bottom: none;
}

.rule-row:hover {
	background: var(--p-surface-100);
}

.rule-name {
	font-weight: 600;
}

.rule-id {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
