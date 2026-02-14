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
		<DecisionDetail v-if="decisionsStore.selectedDecision" :decision="decisionsStore.selectedDecision" />
	</PageContainer>
</template>
