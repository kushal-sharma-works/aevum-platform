<script setup lang="ts">
import { computed } from 'vue'
import type { Decision } from '@/types/decision'

const props = defineProps<{ decisions: ReadonlyArray<Decision> }>()

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
	<table class="decision-table">
		<thead>
			<tr>
				<th>Decision ID</th>
				<th>Rule</th>
				<th>Version</th>
				<th>Status</th>
				<th>Evaluated At</th>
			</tr>
		</thead>
		<tbody>
			<tr v-for="row in rows" :key="row.decisionId">
				<td><RouterLink :to="`/decisions/${row.decisionId}`" class="link">{{ row.decisionId }}</RouterLink></td>
				<td>{{ row.ruleId }}</td>
				<td>{{ row.ruleVersion }}</td>
				<td><span class="p-tag" :class="`p-tag-${statusColor(row.status) === 'positive' ? 'success' : statusColor(row.status) === 'warning' ? 'warn' : 'danger'}`">{{ row.status }}</span></td>
				<td>{{ row.evaluatedAt }}</td>
			</tr>
		</tbody>
	</table>
</template>

<style scoped>
.decision-table {
	width: 100%;
	border-collapse: collapse;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	overflow: hidden;
	background: var(--p-content-background);
}

.decision-table th,
.decision-table td {
	text-align: left;
	padding: 0.625rem;
	border-bottom: 1px solid var(--p-content-border-color);
}

.decision-table tr:last-child td {
	border-bottom: none;
}

.link {
	text-decoration: none;
	color: var(--p-primary-color);
	font-weight: 600;
}
</style>
