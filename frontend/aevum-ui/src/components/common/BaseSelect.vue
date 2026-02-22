<script setup lang="ts">
const model = defineModel<string>({ required: true })
defineProps<{ options: ReadonlyArray<{ readonly label: string; readonly value: string }>; label?: string }>()
</script>

<template>
	<label class="base-field">
		<span v-if="label" class="base-field-label">{{ label }}</span>
		<select v-model="model" class="p-inputtext p-component base-input">
			<option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
		</select>
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
