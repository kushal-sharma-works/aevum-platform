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
		<div class="space-y-2">
			<RouterLink
				v-for="event in eventsStore.events"
				:key="event.eventId"
				:to="`/events/${event.eventId}`"
				class="block rounded border border-slate-800 p-2"
			>
				{{ event.eventType }} — {{ event.eventId }}
			</RouterLink>
		</div>
	</PageContainer>
</template>
