<script setup lang="ts">
defineProps<{ label: string; value: string | number }>()
</script>

<template>
	<div class="rounded border border-slate-800 p-3">
		<div class="text-xs text-slate-400">{{ label }}</div>
		<div class="text-xl font-semibold">{{ value }}</div>
	</div>
</template>
