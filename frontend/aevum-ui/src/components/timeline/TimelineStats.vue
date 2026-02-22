<script setup lang="ts">
import type { TimelineStats } from '@/types/timeline'

defineProps<{ stats: TimelineStats }>()
</script>

<template>
	<div class="timeline-stats">
		<div class="stat-card">Events: {{ stats.eventCount }}</div>
		<div class="stat-card">Decisions: {{ stats.decisionCount }}</div>
		<div class="stat-card">Statuses: {{ Object.keys(stats.statusCounts).length }}</div>
	</div>
</template>

<style scoped>
.timeline-stats {
	display: grid;
	grid-template-columns: repeat(3, minmax(0, 1fr));
	gap: 0.75rem;
}

@media (max-width: 900px) {
	.timeline-stats {
		grid-template-columns: 1fr;
	}
}

.stat-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}
</style>
