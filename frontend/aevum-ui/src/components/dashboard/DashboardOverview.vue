<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'

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
	<section class="q-gutter-md">
		<div class="row q-col-gutter-md">
			<div class="col-12 col-sm-6 col-lg-3">
				<q-card flat bordered>
					<q-card-section>
						<div class="text-caption text-grey-7">Total Events</div>
						<div class="text-h6">{{ totals.events }}</div>
					</q-card-section>
				</q-card>
			</div>
			<div class="col-12 col-sm-6 col-lg-3">
				<q-card flat bordered>
					<q-card-section>
						<div class="text-caption text-grey-7">Total Decisions</div>
						<div class="text-h6">{{ totals.decisions }}</div>
					</q-card-section>
				</q-card>
			</div>
			<div class="col-12 col-sm-6 col-lg-3">
				<q-card flat bordered>
					<q-card-section>
						<div class="text-caption text-grey-7">Active Rules</div>
						<div class="text-h6">{{ totals.activeRules }}</div>
					</q-card-section>
				</q-card>
			</div>
			<div class="col-12 col-sm-6 col-lg-3">
				<q-card flat bordered>
					<q-card-section>
						<div class="text-caption text-grey-7">Avg Latency</div>
						<div class="text-h6">{{ totals.avgLatency }}ms</div>
					</q-card-section>
				</q-card>
			</div>
		</div>

		<div class="row q-col-gutter-md">
			<div class="col-12 col-md-4">
				<q-banner rounded class="bg-blue-1 text-blue-9">IngestionRateChart placeholder</q-banner>
			</div>
			<div class="col-12 col-md-4">
				<q-banner rounded class="bg-teal-1 text-teal-9">DecisionLatencyChart placeholder</q-banner>
			</div>
			<div class="col-12 col-md-4">
				<q-banner rounded class="bg-orange-1 text-orange-9">ErrorRateChart placeholder</q-banner>
			</div>
		</div>
	</section>
</template>
