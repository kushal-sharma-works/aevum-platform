<script setup lang="ts">
defineProps<{ step: { readonly timestamp: string; readonly message: string } }>()
</script>

<template>
	<div class="rounded border border-slate-800 p-2 text-sm">
		<div class="text-xs text-slate-400">{{ step.timestamp }}</div>
		<div>{{ step.message }}</div>
	</div>
</template>
