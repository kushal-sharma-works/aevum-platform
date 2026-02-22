<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import DecisionDetail from '@/components/decisions/DecisionDetail.vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import { useDecisionsStore } from '@/stores/decisions'

const route = useRoute()
const decisionsStore = useDecisionsStore()

onMounted(async () => {
	await decisionsStore.fetchDecision(String(route.params.decisionId))
})
</script>

<template>
	<PageContainer title="Decision Detail">
		<div v-if="decisionsStore.isLoading" class="view-info">Loading decision...</div>
		<div v-else-if="decisionsStore.error" class="view-error">{{ decisionsStore.error }}</div>
		<div v-else-if="!decisionsStore.selectedDecision" class="view-info">No decision found.</div>
		<DecisionDetail v-else :decision="decisionsStore.selectedDecision" />
	</PageContainer>
</template>

<style scoped>
.view-info {
	font-size: 0.9rem;
	color: var(--p-text-muted-color);
}

.view-error {
	padding: 0.55rem 0.7rem;
	border-radius: 0.5rem;
	border: 1px solid color-mix(in oklab, var(--p-red-500) 40%, var(--p-content-border-color));
	background: color-mix(in oklab, var(--p-red-500) 10%, transparent);
	color: var(--p-red-500);
	font-size: 0.85rem;
}
</style>
