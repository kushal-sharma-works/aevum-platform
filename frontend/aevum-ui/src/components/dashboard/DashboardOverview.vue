<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useDecisionsStore } from '@/stores/decisions'
import { useEventsStore } from '@/stores/events'
import { useRulesStore } from '@/stores/rules'

const eventsStore = useEventsStore()
const decisionsStore = useDecisionsStore()
const rulesStore = useRulesStore()

const totals = computed(() => {
	const decisions = decisionsStore.decisions
	const avgLatency = decisions.length
		? Math.round(decisions.reduce((sum, decision) => sum + (decision.trace?.durationMs ?? 0), 0) / decisions.length)
		: 0

	return {
		events: eventsStore.events.length,
		decisions: decisions.length,
		activeRules: rulesStore.rules.length,
		avgLatency
	}
})

onMounted(async () => {
	await Promise.all([
		eventsStore.fetchStreamEvents('default', undefined, 200),
		rulesStore.fetchRules(),
		decisionsStore.fetchDecisions({ page: 1, pageSize: 200 })
	])
})
</script>

<template>
	<section class="dashboard-overview">
		<div class="metrics-grid">
			<section class="overview-card">
				<div class="overview-label">Total Events</div>
				<div class="overview-value">{{ totals.events }}</div>
			</section>
			<section class="overview-card">
				<div class="overview-label">Total Decisions</div>
				<div class="overview-value">{{ totals.decisions }}</div>
			</section>
			<section class="overview-card">
				<div class="overview-label">Active Rules</div>
				<div class="overview-value">{{ totals.activeRules }}</div>
			</section>
			<section class="overview-card">
				<div class="overview-label">Avg Latency</div>
				<div class="overview-value">{{ totals.avgLatency }}ms</div>
			</section>
		</div>

		<div class="charts-grid">
			<section class="chart-placeholder">IngestionRateChart placeholder</section>
			<section class="chart-placeholder">DecisionLatencyChart placeholder</section>
			<section class="chart-placeholder">ErrorRateChart placeholder</section>
		</div>
	</section>
</template>

<style scoped>
.dashboard-overview {
	display: grid;
	gap: 0.85rem;
}

.metrics-grid {
	display: grid;
	grid-template-columns: repeat(4, minmax(0, 1fr));
	gap: 0.75rem;
}

.charts-grid {
	display: grid;
	grid-template-columns: repeat(3, minmax(0, 1fr));
	gap: 0.75rem;
}

@media (max-width: 1100px) {
	.metrics-grid,
	.charts-grid {
		grid-template-columns: 1fr;
	}
}

.overview-card,
.chart-placeholder {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.overview-label {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}

.overview-value {
	font-size: 1.2rem;
	font-weight: 600;
}

.chart-placeholder {
	font-size: 0.9rem;
	color: var(--p-text-muted-color);
}
</style>
