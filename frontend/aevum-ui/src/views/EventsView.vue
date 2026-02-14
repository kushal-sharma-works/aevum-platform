<script setup lang="ts">
import { onMounted } from 'vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import { useEventsStore } from '@/stores/events'

const eventsStore = useEventsStore()

onMounted(async () => {
	await eventsStore.fetchStreamEvents('default')
})
</script>

<template>
	<PageContainer title="Events" description="Stream events with cursor pagination">
		<q-list bordered separator>
			<q-item v-for="event in eventsStore.events" :key="event.eventId" :to="`/events/${event.eventId}`" clickable>
				<q-item-section>
					<q-item-label class="text-weight-medium">{{ event.eventType }}</q-item-label>
					<q-item-label caption>{{ event.eventId }}</q-item-label>
				</q-item-section>
			</q-item>
		</q-list>
	</PageContainer>
</template>
