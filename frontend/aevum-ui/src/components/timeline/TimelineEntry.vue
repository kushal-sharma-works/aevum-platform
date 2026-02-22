<script setup lang="ts">
import { computed } from 'vue'
import type { TimelineEntry } from '@/types/timeline'

const props = defineProps<{ entry: TimelineEntry }>()

const color = computed(() => {
	if (props.entry.kind === 'event') {
		return 'primary'
	}

	switch (props.entry.item.status) {
		case 'Evaluated':
			return 'positive'
		case 'Skipped':
			return 'warning'
		default:
			return 'negative'
	}
})
</script>

<template>
	<div class="timeline-entry">
		<span class="entry-dot" :class="`entry-dot--${String(color)}`" />
		<div>
			<div>{{ entry.timestamp }}</div>
			<div class="entry-kind">{{ entry.kind }}</div>
		</div>
	</div>
</template>

<style scoped>
.timeline-entry {
	display: flex;
	gap: 0.65rem;
	align-items: center;
	padding: 0.45rem 0.6rem;
	border-radius: 0.5rem;
	background: var(--p-content-background);
	border: 1px solid var(--p-content-border-color);
}

.entry-dot {
	width: 0.65rem;
	height: 0.65rem;
	border-radius: 999px;
	display: inline-block;
}

.entry-dot--primary {
	background: var(--p-primary-color);
}

.entry-dot--positive {
	background: var(--p-green-500);
}

.entry-dot--warning {
	background: var(--p-amber-500);
}

.entry-dot--negative {
	background: var(--p-red-500);
}

.entry-kind {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
