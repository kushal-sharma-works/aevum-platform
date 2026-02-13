<script setup lang="ts">
import { useRouter } from 'vue-router'
import PageContainer from '@/components/layout/PageContainer.vue'
import RuleCreateForm from '@/components/rules/RuleCreateForm.vue'
import { useRulesStore } from '@/stores/rules'
import type { RuleCondition } from '@/types/rule'

const router = useRouter()
const rulesStore = useRulesStore()

async function submit(payload: { name: string; description: string; conditions: RuleCondition[] }): Promise<void> {
	const created = await rulesStore.createRule({
		name: payload.name,
		description: payload.description,
		conditions: payload.conditions,
		actions: []
	})
	if (created) {
		await router.push(`/rules/${created.ruleId}`)
	}
}
</script>

<template>
	<PageContainer title="Create Rule">
		<RuleCreateForm @submit="submit" />
	</PageContainer>
</template>
