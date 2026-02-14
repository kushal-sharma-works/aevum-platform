<script setup lang="ts">
import { onMounted } from 'vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import TimelineViewer from '@/components/timeline/TimelineViewer.vue'
import { useTimelineStore } from '@/stores/timeline'

const timelineStore = useTimelineStore()

onMounted(async () => {
	await timelineStore.fetchTimeline()
})
</script>

<template>
	<PageContainer title="Timeline" description="Virtualized timeline of events and decisions">
		<TimelineViewer />
	</PageContainer>
</template>
