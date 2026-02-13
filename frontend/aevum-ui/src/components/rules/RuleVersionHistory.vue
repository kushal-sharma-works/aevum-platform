<script setup lang="ts">
import type { Rule } from '@/types/rule'

defineProps<{ versions: ReadonlyArray<Rule> }>()
</script>

<template>
	<ul class="space-y-1">
		<li v-for="version in versions" :key="`${version.ruleId}-${version.version}`" class="rounded bg-slate-900 px-2 py-1 text-sm">
			v{{ version.version }} - {{ version.isActive ? 'active' : 'inactive' }}
		</li>
	</ul>
</template>
