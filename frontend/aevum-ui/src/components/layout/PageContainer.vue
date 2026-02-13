<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>

<template>
	<section class="space-y-4">
		<header class="space-y-1">
			<h1 class="text-xl font-semibold">{{ title }}</h1>
			<p v-if="description" class="text-sm text-slate-400">{{ description }}</p>
		</header>
		<div>
			<slot />
		</div>
	</section>
</template>
