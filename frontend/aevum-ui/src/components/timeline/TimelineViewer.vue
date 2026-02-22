<script setup lang="ts">
import { computed, ref } from 'vue'
import { useVirtualList } from '@/composables/useVirtualList'
import { useTimelineStore } from '@/stores/timeline'
import TimelineEntry from './TimelineEntry.vue'

const timelineStore = useTimelineStore()
const scrollContainer = ref<HTMLElement | null>(null)

const items = computed(() => [...timelineStore.entries])
const { visibleItems, listStyle, startIndex, itemHeight } = useVirtualList(items, 56, 560)

const onScroll = () => {
	if (!scrollContainer.value) {
		return
	}
}
</script>

<template>
	<section class="timeline-viewer">
		<div class="timeline-header">
			<div class="timeline-title">Timeline Viewer</div>
			<span class="timeline-count">{{ timelineStore.entries.length }} entries</span>
		</div>

		<section class="timeline-card">
			<div ref="scrollContainer" class="timeline-scroll" @scroll="onScroll">
				<div :style="listStyle">
					<div
						v-for="virtual in visibleItems"
						:key="virtual.index"
						class="absolute left-0 right-0 timeline-entry-row"
						:style="{ top: `${virtual.index * itemHeight}px` }"
					>
						<TimelineEntry :entry="virtual.item" />
					</div>
					</div>
				</div>
		</section>

		<div class="timeline-caption">Rendered from index {{ startIndex }} with virtual scrolling.</div>
	</section>
</template>

<style scoped>
.timeline-viewer {
	display: grid;
	gap: 0.75rem;
}

.timeline-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 0.75rem;
}

.timeline-title {
	font-size: 1rem;
	font-weight: 600;
}

.timeline-count {
	display: inline-flex;
	align-items: center;
	padding: 0.2rem 0.5rem;
	border-radius: 999px;
	background: color-mix(in oklab, var(--p-primary-color) 14%, transparent);
	color: var(--p-primary-color);
	font-size: 0.8rem;
	font-weight: 600;
}

.timeline-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
}

.timeline-scroll {
	height: 560px;
	overflow: auto;
}

.timeline-entry-row {
	padding: 0 0.5rem;
}

.timeline-caption {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
