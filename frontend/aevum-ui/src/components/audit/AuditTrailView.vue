<script setup lang="ts">
import type { AuditTrail } from '@/types/audit'
import AuditStep from './AuditStep.vue'

defineProps<{ auditTrail: AuditTrail }>()
</script>

<template>
	<q-list bordered separator>
		<q-item v-for="(step, index) in auditTrail.timeline" :key="index">
			<q-item-section>
				<AuditStep :step="step" />
			</q-item-section>
		</q-item>
	</q-list>
</template>
