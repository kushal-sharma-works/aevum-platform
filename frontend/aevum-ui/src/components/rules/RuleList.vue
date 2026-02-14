<script setup lang="ts">
import type { Rule } from '@/types/rule'

defineProps<{ rules: ReadonlyArray<Rule> }>()
</script>

<template>
	<q-list bordered separator>
		<q-item
			v-for="rule in rules"
			:key="`${rule.ruleId}-${rule.version}`"
			:to="`/rules/${rule.ruleId}`"
			clickable
		>
			<q-item-section>
				<q-item-label class="text-weight-medium">{{ rule.name }}</q-item-label>
				<q-item-label caption>{{ rule.ruleId }}</q-item-label>
			</q-item-section>
			<q-item-section side>
				<q-badge color="primary">v{{ rule.version }}</q-badge>
			</q-item-section>
		</q-item>
	</q-list>
</template>
