<script setup lang="ts">
import { ref } from 'vue'
import type { DecisionTrace } from '@/types/decision'

defineProps<{ trace: DecisionTrace }>()
const expanded = ref<Record<number, boolean>>({})
</script>

<template>
	<div class="q-gutter-sm">
		<div class="text-subtitle2 text-grey-8">Duration: {{ trace.durationMs }}ms</div>
		<q-list bordered separator>
			<q-item v-for="(step, idx) in trace.steps" :key="`${step.conditionField}-${idx}`" class="items-start">
				<q-item-section>
					<div class="row q-col-gutter-sm items-center">
						<div class="col">{{ step.conditionField }}</div>
						<div class="col">{{ step.operator }}</div>
						<div class="col">{{ String(step.expectedValue) }}</div>
						<div class="col">{{ String(step.actualValue) }}</div>
						<div class="col-auto">
							<q-badge :color="step.matched ? 'positive' : 'negative'">{{ step.matched ? '✓' : '✗' }}</q-badge>
						</div>
					</div>
					<q-btn flat dense color="primary" label="Reasoning" @click="expanded[idx] = !expanded[idx]" />
					<div v-if="expanded[idx]" class="text-caption text-grey-7 q-mt-xs">{{ step.reasoning }}</div>
				</q-item-section>
			</q-item>
		</q-list>
	</div>
</template>
