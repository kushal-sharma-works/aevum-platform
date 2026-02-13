<script setup lang="ts">
import type { Decision } from '@/types/decision'
import DecisionCard from './DecisionCard.vue'

defineProps<{ decisions: ReadonlyArray<Decision> }>()
</script>

<template>
	<div class="space-y-2">
		<DecisionCard v-for="decision in decisions" :key="decision.decisionId" :decision="decision" />
	</div>
</template>
