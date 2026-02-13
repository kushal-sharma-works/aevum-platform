<script setup lang="ts">
import type { AuditTrail } from '@/types/audit'
import AuditStep from './AuditStep.vue'

defineProps<{ auditTrail: AuditTrail }>()
</script>

<template>
	<div class="space-y-2">
		<AuditStep v-for="(step, index) in auditTrail.timeline" :key="index" :step="step" />
	</div>
</template>
