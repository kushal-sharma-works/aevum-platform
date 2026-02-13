<script setup lang="ts">
import type { Rule } from '@/types/rule'

defineProps<{ rules: ReadonlyArray<Rule> }>()
</script>

<template>
	<div class="space-y-2">
		<RouterLink
			v-for="rule in rules"
			:key="`${rule.ruleId}-${rule.version}`"
			:to="`/rules/${rule.ruleId}`"
			class="block rounded border border-slate-800 p-3"
		>
			<div class="flex justify-between">
				<span>{{ rule.name }}</span>
				<span class="text-xs text-slate-400">v{{ rule.version }}</span>
			</div>
		</RouterLink>
	</div>
</template>
