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
	<section class="q-gutter-sm">
		<div class="row items-center justify-between">
			<div class="text-subtitle1 text-weight-medium">Timeline Viewer</div>
			<q-badge color="primary">{{ timelineStore.entries.length }} entries</q-badge>
		</div>

		<q-card flat bordered>
			<q-card-section class="q-pa-none">
				<div ref="scrollContainer" class="timeline-scroll" @scroll="onScroll">
					<div :style="listStyle">
						<div
							v-for="virtual in visibleItems"
							:key="virtual.index"
							class="absolute left-0 right-0 q-px-sm"
							:style="{ top: `${virtual.index * itemHeight}px` }"
						>
							<TimelineEntry :entry="virtual.item" />
						</div>
					</div>
				</div>
			</q-card-section>
		</q-card>

		<div class="text-caption text-grey-7">Rendered from index {{ startIndex }} with virtual scrolling.</div>
	</section>
</template>

<style scoped>
.timeline-scroll {
	height: 560px;
	overflow: auto;
}
</style>
