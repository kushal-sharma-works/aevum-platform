<script setup lang="ts">
import { onMounted } from 'vue'
import BaseEmptyState from '@/components/common/BaseEmptyState.vue'
import DecisionList from '@/components/decisions/DecisionList.vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import { useDecisionsStore } from '@/stores/decisions'

const decisionsStore = useDecisionsStore()

onMounted(async () => {
	await decisionsStore.fetchDecisions({ page: 1, pageSize: 20 })
})
</script>

<template>
	<PageContainer title="Decisions" description="Paginated decision inspector">
		<div v-if="decisionsStore.isLoading" class="view-info">Loading decisions...</div>
		<div v-else-if="decisionsStore.error" class="view-error">{{ decisionsStore.error }}</div>
		<BaseEmptyState
			v-else-if="decisionsStore.decisions.length === 0"
			title="No decisions"
			description="No decision records are available yet."
		/>
		<DecisionList v-else :decisions="decisionsStore.decisions" />
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
