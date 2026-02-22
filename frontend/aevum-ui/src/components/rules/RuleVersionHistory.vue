<script setup lang="ts">
import type { Rule } from '@/types/rule'

defineProps<{ versions: ReadonlyArray<Rule> }>()
</script>

<template>
	<ul class="version-list">
		<li v-for="version in versions" :key="`${version.ruleId}-${version.version}`" class="version-item">
			<span>v{{ version.version }}</span>
			<span class="version-pill" :class="version.isActive ? 'version-pill--active' : ''">
				{{ version.isActive ? 'active' : 'inactive' }}
			</span>
		</li>
	</ul>
</template>

<style scoped>
.version-list {
	list-style: none;
	margin: 0;
	padding: 0;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
}

.version-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 0.65rem 0.75rem;
	border-bottom: 1px solid var(--p-content-border-color);
}

.version-item:last-child {
	border-bottom: 0;
}

.version-pill {
	border-radius: 999px;
	padding: 0.15rem 0.5rem;
	font-size: 0.75rem;
	border: 1px solid var(--p-content-border-color);
	color: var(--p-text-muted-color);
}

.version-pill--active {
	color: var(--p-green-500);
	border-color: color-mix(in oklab, var(--p-green-500) 45%, var(--p-content-border-color));
	background: color-mix(in oklab, var(--p-green-500) 11%, transparent);
}
</style>
