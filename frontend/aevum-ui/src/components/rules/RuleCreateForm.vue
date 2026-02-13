<script setup lang="ts">
import { ref } from 'vue'
import BaseButton from '@/components/common/BaseButton.vue'
import ConditionBuilder from './ConditionBuilder.vue'
import type { RuleCondition } from '@/types/rule'

const emit = defineEmits<{ (e: 'submit', value: { name: string; description: string; conditions: RuleCondition[] }): void }>()
const name = ref('')
const description = ref('')
const conditions = ref<RuleCondition[]>([])

function submit(): void {
	emit('submit', { name: name.value, description: description.value, conditions: conditions.value })
}
</script>

<template>
	<div class="space-y-3 rounded border border-slate-800 p-3">
		<input v-model="name" class="w-full rounded bg-slate-900 px-2 py-1" placeholder="Rule name" />
		<input v-model="description" class="w-full rounded bg-slate-900 px-2 py-1" placeholder="Description" />
		<ConditionBuilder v-model="conditions" />
		<BaseButton @click="submit">Create Rule</BaseButton>
	</div>
</template>
