<script setup lang="ts">
import { ref } from 'vue'
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
	<q-card flat bordered>
		<q-card-section class="q-gutter-sm">
			<q-input v-model="name" label="Rule name" outlined />
			<q-input v-model="description" label="Description" outlined />
			<ConditionBuilder v-model="conditions" />
			<q-btn color="primary" label="Create Rule" @click="submit" />
		</q-card-section>
	</q-card>
</template>
