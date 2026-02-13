<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import { useRulesStore } from '@/stores/rules'

const route = useRoute()
const rulesStore = useRulesStore()

onMounted(async () => {
	await rulesStore.fetchRule(String(route.params.ruleId))
	await rulesStore.fetchRuleVersions(String(route.params.ruleId))
})
</script>

<template>
	<PageContainer title="Rule Detail">
		<BaseJsonViewer v-if="rulesStore.selectedRule" :value="rulesStore.selectedRule" />
		<BaseJsonViewer
			v-if="rulesStore.ruleVersions[String(route.params.ruleId)]"
			:value="rulesStore.ruleVersions[String(route.params.ruleId)]"
			label="Versions"
		/>
	</PageContainer>
</template>
