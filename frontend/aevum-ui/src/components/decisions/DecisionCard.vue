<script setup lang="ts">
import type { Decision } from '@/types/decision'

defineProps<{ decision: Decision }>()
</script>

<template>
	<RouterLink
		:to="`/decisions/${decision.decisionId}`"
		class="block rounded border border-slate-800 p-3 hover:border-slate-700"
	>
		<div class="flex items-center justify-between">
			<span class="font-medium">{{ decision.decisionId }}</span>
			<span class="text-xs text-slate-400">{{ decision.status }}</span>
		</div>
		<p class="mt-1 text-xs text-slate-400">Rule {{ decision.ruleId }} v{{ decision.ruleVersion }}</p>
	</RouterLink>
</template>
