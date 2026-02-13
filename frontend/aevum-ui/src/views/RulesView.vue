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
		<div class="mb-2">
			<RouterLink to="/rules/create" class="text-sm text-blue-300">Create Rule</RouterLink>
		</div>
		<RuleList :rules="rulesStore.rules" />
	</PageContainer>
</template>
