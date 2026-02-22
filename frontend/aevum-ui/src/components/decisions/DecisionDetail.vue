<script setup lang="ts">
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'
import type { Decision } from '@/types/decision'
import DecisionTrace from './DecisionTrace.vue'

defineProps<{ decision: Decision }>()
</script>

<template>
	<div class="decision-detail">
		<section class="detail-card">
			<div class="detail-title">Decision {{ decision.decisionId }}</div>
			<div class="detail-meta">Rule {{ decision.ruleId }} v{{ decision.ruleVersion }}</div>
		</section>
		<div class="detail-grid">
			<div>
				<BaseJsonViewer :value="decision.input" label="Input" />
			</div>
			<div>
				<BaseJsonViewer :value="decision.output" label="Output" />
			</div>
		</div>
		<DecisionTrace :trace="decision.trace" />
	</div>
</template>

<style scoped>
.decision-detail {
	display: grid;
	gap: 0.85rem;
}

.detail-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.detail-title {
	font-size: 1rem;
	font-weight: 600;
}

.detail-meta {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
	margin-top: 0.2rem;
}

.detail-grid {
	display: grid;
	grid-template-columns: repeat(2, minmax(0, 1fr));
	gap: 0.75rem;
}

@media (max-width: 1024px) {
	.detail-grid {
		grid-template-columns: 1fr;
	}
}
</style>
