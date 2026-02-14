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
	<div class="space-y-2">
		<div v-for="(condition, index) in model" :key="index" class="grid grid-cols-12 gap-2">
			<input v-model="condition.field" class="col-span-4 rounded bg-slate-900 px-2 py-1" placeholder="Field" />
			<select v-model="condition.operator" class="col-span-3 rounded bg-slate-900 px-2 py-1">
				<option v-for="op in operators" :key="op" :value="op">{{ op }}</option>
			</select>
			<input v-model="condition.value" class="col-span-4 rounded bg-slate-900 px-2 py-1" placeholder="Value" />
			<button class="col-span-1 rounded bg-rose-700 px-2" @click="remove(index)">-</button>
		</div>
		<BaseButton variant="secondary" @click="add">Add Condition</BaseButton>
	</div>
</template>
