<script setup lang="ts">
import type { Rule } from '@/types/rule'

defineProps<{ rule: Rule }>()
</script>

<template>
	<div class="rounded border border-slate-800 p-3">
		<h3 class="font-semibold">{{ rule.name }} v{{ rule.version }}</h3>
		<p class="mt-1 text-sm text-slate-400">{{ rule.description }}</p>
	</div>
</template>
