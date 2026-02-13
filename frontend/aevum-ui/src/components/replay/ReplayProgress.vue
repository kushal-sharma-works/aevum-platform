<script setup lang="ts">
defineProps<{ progress: number }>()
</script>

<template>
	<div>
		<div class="mb-1 text-xs text-slate-400">Progress {{ progress.toFixed(0) }}%</div>
		<div class="h-2 rounded bg-slate-800">
			<div class="h-2 rounded bg-blue-600" :style="{ width: `${progress}%` }" />
		</div>
	</div>
</template>
