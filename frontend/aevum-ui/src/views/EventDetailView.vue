<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import { useEventsStore } from '@/stores/events'

const route = useRoute()
const eventsStore = useEventsStore()

onMounted(async () => {
	const id = String(route.params.eventId)
	await eventsStore.fetchEvent(id)
})
</script>

<template>
	<PageContainer title="Event Detail">
		<div v-if="eventsStore.isLoading" class="view-info">Loading event detail...</div>
		<div v-else-if="eventsStore.error" class="view-error">{{ eventsStore.error }}</div>
		<div v-else-if="!eventsStore.selectedEvent" class="view-info">No event found for this identifier.</div>
		<BaseJsonViewer v-else :value="eventsStore.selectedEvent" />
	</PageContainer>
</template>

<style scoped>
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
</style>
