<script setup lang="ts">
import type { Decision } from '@/types/decision'
import type { Event } from '@/types/event'

defineProps<{ items: ReadonlyArray<Event | Decision> }>()
</script>

<template>
	<ul class="feed-list">
		<li v-for="(item, index) in items" :key="index" class="feed-item">
			{{ 'decisionId' in item ? `Decision ${item.decisionId}` : `Event ${item.eventId}` }}
		</li>
	</ul>
</template>

<style scoped>
.feed-list {
	list-style: none;
	margin: 0;
	padding: 0;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	max-height: 240px;
	overflow: auto;
}

.feed-item {
	padding: 0.65rem 0.8rem;
	border-bottom: 1px solid var(--p-content-border-color);
}

.feed-item:last-child {
	border-bottom: 0;
}
</style>
