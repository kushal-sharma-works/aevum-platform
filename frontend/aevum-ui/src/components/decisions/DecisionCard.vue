<script setup lang="ts">
import type { Decision } from '@/types/decision'

defineProps<{ decision: Decision }>()
</script>

<template>
	<RouterLink :to="`/decisions/${decision.decisionId}`" class="block">
		<section class="decision-card">
			<div class="decision-head">
				<span class="decision-id">{{ decision.decisionId }}</span>
				<span class="status-pill">{{ decision.status }}</span>
			</div>
			<div class="decision-meta">Rule {{ decision.ruleId }} v{{ decision.ruleVersion }}</div>
		</section>
	</RouterLink>
</template>

<style scoped>
.decision-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.decision-head {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 0.5rem;
}

.decision-id {
	font-weight: 600;
}

.status-pill {
	border: 1px solid var(--p-content-border-color);
	border-radius: 999px;
	padding: 0.15rem 0.5rem;
	font-size: 0.75rem;
	color: var(--p-text-muted-color);
}

.decision-meta {
	margin-top: 0.3rem;
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
