<script setup lang="ts">
import { onMounted } from 'vue'
import BaseEmptyState from '@/components/common/BaseEmptyState.vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import { useEventsStore } from '@/stores/events'

const eventsStore = useEventsStore()

onMounted(async () => {
	await eventsStore.fetchStreamEvents('default')
})
</script>

<template>
	<PageContainer title="Events" description="Stream events with cursor pagination">
		<div v-if="eventsStore.isLoading" class="view-info">Loading events...</div>
		<div v-else-if="eventsStore.error" class="view-error">{{ eventsStore.error }}</div>
		<BaseEmptyState
			v-else-if="eventsStore.events.length === 0"
			title="No events"
			description='No events found for stream "default".'
		/>
		<div v-else class="events-list">
			<RouterLink v-for="event in eventsStore.events" :key="event.eventId" :to="`/events/${event.eventId}`" class="event-row">
				<div class="event-type">{{ event.eventType }}</div>
				<div class="event-id">{{ event.eventId }}</div>
			</RouterLink>
		</div>
	</PageContainer>
</template>

<style scoped>
.events-list {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	overflow: hidden;
}

.view-info {
	font-size: 0.9rem;
	color: var(--p-text-muted-color);
}

.view-error {
	padding: 0.55rem 0.7rem;
	border-radius: 0.5rem;
	border: 1px solid color-mix(in oklab, var(--p-red-500) 40%, var(--p-content-border-color));
	background: color-mix(in oklab, var(--p-red-500) 10%, transparent);
	color: var(--p-red-500);
	font-size: 0.85rem;
}

.event-row {
	display: block;
	padding: 0.75rem;
	text-decoration: none;
	color: inherit;
	border-bottom: 1px solid var(--p-content-border-color);
}

.event-row:last-child {
	border-bottom: none;
}

.event-row:hover {
	background: var(--p-surface-100);
}

.event-type {
	font-weight: 600;
}

.event-id {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
