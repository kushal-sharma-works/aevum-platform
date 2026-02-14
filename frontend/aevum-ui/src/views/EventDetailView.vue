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
		<BaseJsonViewer v-if="eventsStore.selectedEvent" :value="eventsStore.selectedEvent" />
	</PageContainer>
</template>
