<script setup lang="ts">
const model = defineModel<string>({ required: true })
defineProps<{ label?: string; placeholder?: string; error?: string }>()
</script>

<template>
	<label class="block space-y-1">
		<span v-if="label" class="text-xs text-slate-300">{{ label }}</span>
		<input v-model="model" :placeholder="placeholder" class="w-full rounded bg-slate-900 px-3 py-2 text-sm" />
		<span v-if="error" class="text-xs text-rose-300">{{ error }}</span>
	</label>
</template>
