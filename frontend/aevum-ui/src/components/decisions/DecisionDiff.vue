<script setup lang="ts">
import { computed } from 'vue'
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'
import { diffJson } from '@/utils/jsonDiff'

const props = defineProps<{ left: Record<string, unknown>; right: Record<string, unknown> }>()
const diffs = computed(() => diffJson(props.left, props.right))
</script>

<template>
	<section class="decision-diff">
		<div class="diff-grid">
			<section class="diff-card">
				<div class="diff-title">T1 / V1</div>
				<BaseJsonViewer :value="left" />
			</section>
			<section class="diff-card">
				<div class="diff-title">T2 / V2</div>
				<BaseJsonViewer :value="right" />
			</section>
		</div>
		<div>
			<div class="diff-title">Changes</div>
			<ul class="change-list">
				<li v-for="entry in diffs" :key="entry.path" class="change-item" :class="`change-item--${entry.type}`">
					{{ entry.type }} @ {{ entry.path }}
				</li>
			</ul>
		</div>
	</section>
</template>

<style scoped>
.decision-diff {
	display: grid;
	gap: 0.85rem;
}

.diff-grid {
	display: grid;
	grid-template-columns: repeat(2, minmax(0, 1fr));
	gap: 0.75rem;
}

@media (max-width: 1024px) {
	.diff-grid {
		grid-template-columns: 1fr;
	}
}

.diff-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.diff-title {
	font-weight: 600;
	margin-bottom: 0.45rem;
}

.change-list {
	list-style: none;
	margin: 0;
	padding: 0;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	overflow: hidden;
}

.change-item {
	padding: 0.65rem;
	border-bottom: 1px solid var(--p-content-border-color);
}

.change-item:last-child {
	border-bottom: 0;
}

.change-item--added {
	background: color-mix(in oklab, var(--p-green-500) 12%, transparent);
}

.change-item--removed {
	background: color-mix(in oklab, var(--p-red-500) 12%, transparent);
}

.change-item--changed {
	background: color-mix(in oklab, var(--p-amber-500) 15%, transparent);
}
</style>
