<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import PageContainer from '@/components/layout/PageContainer.vue'
import RuleCreateForm from '@/components/rules/RuleCreateForm.vue'
import { useRulesStore } from '@/stores/rules'
import type { RuleCondition } from '@/types/rule'

const router = useRouter()
const rulesStore = useRulesStore()
const localError = ref<string | null>(null)

async function submit(payload: { name: string; description: string; conditions: RuleCondition[] }): Promise<void> {
	localError.value = null
	if (!payload.name.trim()) {
		localError.value = 'Rule name is required'
		return
	}

	if (payload.conditions.length === 0) {
		localError.value = 'Add at least one condition before creating a rule'
		return
	}

	if (payload.conditions.some((condition) => !condition.field.trim())) {
		localError.value = 'Each condition needs a field'
		return
	}

	const created = await rulesStore.createRule({
		name: payload.name,
		description: payload.description,
		conditions: payload.conditions,
		actions: [{ actionType: 'Approve', parameters: {} }]
	})
	if (created) {
		await router.push(`/rules/${created.ruleId}`)
		return
	}

	localError.value = rulesStore.error ?? 'Failed to create rule'
}
</script>

<template>
	<PageContainer title="Create Rule">
		<div v-if="localError || rulesStore.error" class="create-rule-error">{{ localError ?? rulesStore.error }}</div>
		<RuleCreateForm :submitting="rulesStore.isLoading" @submit="submit" />
	</PageContainer>
</template>

<style scoped>
.create-rule-error {
	margin-bottom: 0.75rem;
	padding: 0.55rem 0.7rem;
	border-radius: 0.5rem;
	border: 1px solid color-mix(in oklab, var(--p-red-500) 40%, var(--p-content-border-color));
	background: color-mix(in oklab, var(--p-red-500) 10%, transparent);
	color: var(--p-red-500);
	font-size: 0.85rem;
}
</style>
