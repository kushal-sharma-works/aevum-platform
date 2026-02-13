<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import MetricCard from './MetricCard.vue'

const totals = ref({
	events: 0,
	decisions: 0,
	activeRules: 0,
	avgLatency: 0
})

let timer: number | undefined

const refresh = () => {
	totals.value = {
		events: totals.value.events + 10,
		decisions: totals.value.decisions + 4,
		activeRules: Math.max(1, totals.value.activeRules || 3),
		avgLatency: 42
	}
}

onMounted(() => {
	refresh()
	timer = window.setInterval(refresh, 30_000)
})

onUnmounted(() => {
	if (timer !== undefined) {
		window.clearInterval(timer)
	}
})
</script>

<template>
	<section class="space-y-4">
		<div class="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
			<MetricCard label="Total Events" :value="totals.events" />
			<MetricCard label="Total Decisions" :value="totals.decisions" />
			<MetricCard label="Active Rules" :value="totals.activeRules" />
			<MetricCard label="Avg Latency" :value="`${totals.avgLatency}ms`" />
		</div>
		<div class="grid gap-3 md:grid-cols-3">
			<div class="rounded border border-slate-800 p-3 text-sm">IngestionRateChart placeholder</div>
			<div class="rounded border border-slate-800 p-3 text-sm">DecisionLatencyChart placeholder</div>
			<div class="rounded border border-slate-800 p-3 text-sm">ErrorRateChart placeholder</div>
		</div>
	</section>
</template>
