<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { queryAuditApi } from '@/api/queryAudit'
import AuditTrailView from '@/components/audit/AuditTrailView.vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import type { AuditTrail } from '@/types/audit'

const route = useRoute()
const auditTrail = ref<AuditTrail | null>(null)

onMounted(async () => {
	auditTrail.value = await queryAuditApi.getAuditTrail(String(route.params.decisionId))
})
</script>

<template>
	<PageContainer title="Audit Trail">
		<AuditTrailView v-if="auditTrail" :audit-trail="auditTrail" />
	</PageContainer>
</template>
