<script setup lang="ts">
import { onMounted } from 'vue'
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
		<DecisionList :decisions="decisionsStore.decisions" />
	</PageContainer>
</template>
