<script setup lang="ts">
import type { AuditTrail } from '@/types/audit'
import AuditStep from './AuditStep.vue'

defineProps<{ auditTrail: AuditTrail }>()
</script>

<template>
	<ul class="audit-list">
		<li v-for="(step, index) in auditTrail.timeline" :key="index" class="audit-item">
			<AuditStep :step="step" />
		</li>
	</ul>
</template>

<style scoped>
.audit-list {
	list-style: none;
	margin: 0;
	padding: 0;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
}

.audit-item {
	padding: 0.65rem;
	border-bottom: 1px solid var(--p-content-border-color);
}

.audit-item:last-child {
	border-bottom: 0;
}
</style>
