<script setup lang="ts">
const model = defineModel<string>({ required: true })
defineProps<{ label?: string; placeholder?: string; error?: string; type?: string }>()
</script>

<template>
	<label class="base-field">
		<span v-if="label" class="base-field-label">{{ label }}</span>
		<input v-model="model" :type="type ?? 'text'" :placeholder="placeholder" class="p-inputtext p-component base-input" />
		<small v-if="error" class="p-error">{{ error }}</small>
	</label>
</template>

<style scoped>
.base-field {
	display: flex;
	flex-direction: column;
	gap: 0.35rem;
}

.base-field-label {
	font-size: 0.85rem;
	color: var(--p-text-muted-color);
}

.base-input {
	width: 100%;
}
</style>
