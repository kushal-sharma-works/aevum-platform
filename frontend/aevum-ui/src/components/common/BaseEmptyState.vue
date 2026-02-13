<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>

<template>
	<div class="rounded border border-dashed border-slate-700 p-6 text-center">
		<h3 class="font-semibold">{{ title }}</h3>
		<p v-if="description" class="mt-1 text-sm text-slate-400">{{ description }}</p>
	</div>
</template>
