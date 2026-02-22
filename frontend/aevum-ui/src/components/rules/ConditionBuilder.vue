<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'
import type { ConditionOperator, RuleCondition } from '@/types/rule'

const model = defineModel<RuleCondition[]>({ required: true })
const emit = defineEmits<{ (e: 'validated', value: RuleCondition[]): void }>()

const operators: ReadonlyArray<ConditionOperator> = [
	'Eq',
	'NotEq',
	'Gt',
	'Gte',
	'Lt',
	'Lte',
	'Contains',
	'NotContains',
	'Regex',
	'In',
	'NotIn'
]

function add(): void {
	model.value = [...model.value, { field: '', operator: 'Eq', value: '' }]
	emit('validated', model.value)
}

function remove(index: number): void {
	model.value = model.value.filter((_, idx) => idx !== index)
	emit('validated', model.value)
}
</script>

<template>
	<div class="condition-builder">
		<div v-for="(condition, index) in model" :key="index" class="condition-row">
			<label class="condition-field">
				<span>Field</span>
				<input v-model="condition.field" class="p-inputtext p-component" />
			</label>
			<label class="condition-field">
				<span>Operator</span>
				<select v-model="condition.operator" class="p-inputtext p-component">
					<option v-for="operator in operators" :key="operator" :value="operator">{{ operator }}</option>
				</select>
			</label>
			<label class="condition-field">
				<span>Value</span>
				<input v-model="condition.value" class="p-inputtext p-component" />
			</label>
			<button type="button" class="p-button p-component p-button-danger p-button-sm" @click="remove(index)">Remove</button>
		</div>
		<BaseButton variant="secondary" @click="add">Add Condition</BaseButton>
	</div>
</template>

<style scoped>
.condition-builder {
	display: grid;
	gap: 0.75rem;
}

.condition-row {
	display: grid;
	grid-template-columns: repeat(4, minmax(0, 1fr));
	gap: 0.75rem;
	align-items: end;
}

@media (max-width: 1100px) {
	.condition-row {
		grid-template-columns: 1fr;
	}
}

.condition-field {
	display: grid;
	gap: 0.35rem;
}

.condition-field span {
	font-size: 0.85rem;
	color: var(--p-text-muted-color);
}
</style>
