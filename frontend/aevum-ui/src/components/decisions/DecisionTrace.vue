<script setup lang="ts">
import { ref } from 'vue'
import type { DecisionTrace } from '@/types/decision'

defineProps<{ trace: DecisionTrace }>()
const expanded = ref<Record<number, boolean>>({})
</script>

<template>
	<div class="decision-trace">
		<div class="trace-duration">Duration: {{ trace.durationMs }}ms</div>
		<ul class="trace-list">
			<li v-for="(step, idx) in trace.steps" :key="`${step.conditionField}-${idx}`" class="trace-item">
				<div class="trace-row">
					<div>{{ step.conditionField }}</div>
					<div>{{ step.operator }}</div>
					<div>{{ String(step.expectedValue) }}</div>
					<div>{{ String(step.actualValue) }}</div>
					<div class="trace-result" :class="step.matched ? 'trace-result--ok' : 'trace-result--bad'">
						{{ step.matched ? '✓' : '✗' }}
					</div>
				</div>
				<button type="button" class="p-button p-component p-button-text" @click="expanded[idx] = !expanded[idx]">
					Reasoning
				</button>
				<div v-if="expanded[idx]" class="trace-reason">{{ step.reasoning }}</div>
			</li>
		</ul>
	</div>
</template>

<style scoped>
.decision-trace {
	display: grid;
	gap: 0.55rem;
}

.trace-duration {
	color: var(--p-text-muted-color);
	font-size: 0.9rem;
}

.trace-list {
	list-style: none;
	margin: 0;
	padding: 0;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
}

.trace-item {
	padding: 0.65rem;
	border-bottom: 1px solid var(--p-content-border-color);
}

.trace-item:last-child {
	border-bottom: 0;
}

.trace-row {
	display: grid;
	grid-template-columns: repeat(5, minmax(0, 1fr));
	gap: 0.5rem;
	align-items: center;
}

.trace-result {
	text-align: center;
	font-weight: 700;
}

.trace-result--ok {
	color: var(--p-green-500);
}

.trace-result--bad {
	color: var(--p-red-500);
}

.trace-reason {
	margin-top: 0.25rem;
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
