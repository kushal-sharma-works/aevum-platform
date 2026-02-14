<script setup lang="ts">
const model = defineModel<string>({ required: true })
defineProps<{ options: ReadonlyArray<{ readonly label: string; readonly value: string }>; label?: string }>()
</script>

<template>
	<label class="block space-y-1">
		<span v-if="label" class="text-xs text-slate-300">{{ label }}</span>
		<select v-model="model" class="w-full rounded bg-slate-900 px-3 py-2 text-sm">
			<option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
		</select>
	</label>
</template>
