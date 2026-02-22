<script setup lang="ts">
import { ref } from 'vue'
import ConditionBuilder from './ConditionBuilder.vue'
import type { RuleCondition } from '@/types/rule'

defineProps<{ submitting?: boolean }>()
const emit = defineEmits<{ (e: 'submit', value: { name: string; description: string; conditions: RuleCondition[] }): void }>()
const name = ref('')
const description = ref('')
const conditions = ref<RuleCondition[]>([])

function submit(): void {
	emit('submit', { name: name.value, description: description.value, conditions: conditions.value })
}
</script>

<template>
	<section class="rule-form-card">
		<label class="rule-field">
			<span>Rule name</span>
			<input v-model="name" class="p-inputtext p-component" />
		</label>
		<label class="rule-field">
			<span>Description</span>
			<input v-model="description" class="p-inputtext p-component" />
		</label>
		<ConditionBuilder v-model="conditions" />
		<button type="button" class="p-button p-component" :disabled="submitting" @click="submit">Create Rule</button>
	</section>
</template>

<style scoped>
.rule-form-card {
	display: grid;
	gap: 0.75rem;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.rule-field {
	display: grid;
	gap: 0.35rem;
}

.rule-field span {
	font-size: 0.85rem;
	color: var(--p-text-muted-color);
}
</style>
