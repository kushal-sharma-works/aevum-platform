<script setup lang="ts">
import { onMounted } from 'vue'
import BaseEmptyState from '@/components/common/BaseEmptyState.vue'
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
		<div class="create-wrap">
			<RouterLink to="/rules/create" class="p-button p-component">Create Rule</RouterLink>
		</div>
		<div v-if="rulesStore.isLoading" class="view-info">Loading rules...</div>
		<div v-else-if="rulesStore.error" class="view-error">{{ rulesStore.error }}</div>
		<BaseEmptyState
			v-else-if="rulesStore.rules.length === 0"
			title="No rules"
			description="No active rules are available."
		/>
		<RuleList v-else :rules="rulesStore.rules" />
	</PageContainer>
</template>

<style scoped>
.create-wrap {
	margin-bottom: 0.6rem;
}

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
