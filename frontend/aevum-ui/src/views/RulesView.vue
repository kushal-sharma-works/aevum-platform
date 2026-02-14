<script setup lang="ts">
import { onMounted } from 'vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import RuleList from '@/components/rules/RuleList.vue'
import { useRulesStore } from '@/stores/rules'

const rulesStore = useRulesStore()

onMounted(async () => {
	await rulesStore.fetchRules()
})
</script>

<template>
	<PageContainer title="Rules" description="Rule version browser">
		<div class="q-mb-sm">
			<q-btn color="primary" icon="add" label="Create Rule" to="/rules/create" />
		</div>
		<RuleList :rules="rulesStore.rules" />
	</PageContainer>
</template>
