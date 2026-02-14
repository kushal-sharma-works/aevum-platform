<script setup lang="ts">
import type { TimelineStats } from '@/types/timeline'

defineProps<{ stats: TimelineStats }>()
</script>

<template>
	<div class="grid gap-2 md:grid-cols-3">
		<div class="rounded border border-slate-800 p-2 text-sm">Events: {{ stats.eventCount }}</div>
		<div class="rounded border border-slate-800 p-2 text-sm">Decisions: {{ stats.decisionCount }}</div>
		<div class="rounded border border-slate-800 p-2 text-sm">Statuses: {{ Object.keys(stats.statusCounts).length }}</div>
	</div>
</template>
