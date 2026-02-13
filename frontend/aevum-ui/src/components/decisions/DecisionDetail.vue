<script setup lang="ts">
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'
import type { Decision } from '@/types/decision'
import DecisionTrace from './DecisionTrace.vue'

defineProps<{ decision: Decision }>()
</script>

<template>
	<div class="space-y-4">
		<div class="rounded border border-slate-800 p-3">
			<h3 class="mb-2 font-semibold">Decision {{ decision.decisionId }}</h3>
			<p class="text-xs text-slate-400">Rule {{ decision.ruleId }} v{{ decision.ruleVersion }}</p>
		</div>
		<div class="grid gap-4 lg:grid-cols-2">
			<BaseJsonViewer :value="decision.input" label="Input" />
			<BaseJsonViewer :value="decision.output" label="Output" />
		</div>
		<DecisionTrace :trace="decision.trace" />
	</div>
</template>
