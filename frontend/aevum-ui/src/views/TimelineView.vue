<script setup lang="ts">
import { onMounted } from 'vue'
import BaseEmptyState from '@/components/common/BaseEmptyState.vue'
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
		<div v-if="timelineStore.isLoading" class="view-info">Loading timeline...</div>
		<div v-else-if="timelineStore.error" class="view-error">{{ timelineStore.error }}</div>
		<BaseEmptyState
			v-else-if="timelineStore.entries.length === 0"
			title="No timeline entries"
			description="No events or decisions found for the selected range."
		/>
		<TimelineViewer v-else />
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
