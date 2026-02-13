<script setup lang="ts">
import { ref } from 'vue'
import type { DecisionTrace } from '@/types/decision'

defineProps<{ trace: DecisionTrace }>()
const expanded = ref<Record<number, boolean>>({})
</script>

<template>
	<div class="space-y-2">
		<div class="text-sm text-slate-300">Duration: {{ trace.durationMs }}ms</div>
		<div class="space-y-1">
			<div
				v-for="(step, idx) in trace.steps"
				:key="`${step.conditionField}-${idx}`"
				:class="[
					'rounded border p-2 text-xs',
					step.matched ? 'border-emerald-700 bg-emerald-950/20' : 'border-rose-700 bg-rose-950/20'
				]"
			>
				<div class="grid grid-cols-5 gap-2">
					<span>{{ step.conditionField }}</span>
					<span>{{ step.operator }}</span>
					<span>{{ String(step.expectedValue) }}</span>
					<span>{{ String(step.actualValue) }}</span>
					<span>{{ step.matched ? '✓' : '✗' }}</span>
				</div>
				<button class="mt-1 text-blue-300" @click="expanded[idx] = !expanded[idx]">Reasoning</button>
				<p v-if="expanded[idx]" class="mt-1 text-slate-300">{{ step.reasoning }}</p>
			</div>
		</div>
	</div>
</template>
