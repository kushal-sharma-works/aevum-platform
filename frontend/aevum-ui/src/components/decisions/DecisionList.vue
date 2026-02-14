<script setup lang="ts">
import { computed } from 'vue'
import type { Decision } from '@/types/decision'

const props = defineProps<{ decisions: ReadonlyArray<Decision> }>()

const columns = [
	{ name: 'decisionId', label: 'Decision ID', field: 'decisionId', align: 'left' as const, sortable: true },
	{ name: 'rule', label: 'Rule', field: 'ruleId', align: 'left' as const, sortable: true },
	{ name: 'version', label: 'Version', field: 'ruleVersion', align: 'left' as const, sortable: true },
	{ name: 'status', label: 'Status', field: 'status', align: 'left' as const, sortable: true },
	{ name: 'evaluatedAt', label: 'Evaluated At', field: 'evaluatedAt', align: 'left' as const, sortable: true }
]

const rows = computed(() => props.decisions)

const statusColor = (status: Decision['status']) => {
	switch (status) {
		case 'Evaluated':
			return 'positive'
		case 'Skipped':
			return 'warning'
		default:
			return 'negative'
	}
}
</script>

<template>
	<q-table
		flat
		bordered
		row-key="decisionId"
		:rows="rows"
		:columns="columns"
		:pagination="{ rowsPerPage: 20 }"
	>
		<template #body-cell-decisionId="scope">
			<q-td :props="scope">
				<RouterLink :to="`/decisions/${scope.row.decisionId}`" class="text-primary text-weight-medium">
					{{ scope.row.decisionId }}
				</RouterLink>
			</q-td>
		</template>

		<template #body-cell-rule="scope">
			<q-td :props="scope">{{ scope.row.ruleId }}</q-td>
		</template>

		<template #body-cell-status="scope">
			<q-td :props="scope">
				<q-badge :color="statusColor(scope.row.status)" text-color="white">{{ scope.row.status }}</q-badge>
			</q-td>
		</template>
	</q-table>
</template>
