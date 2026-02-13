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
	<section class="space-y-2">
		<header class="flex items-center justify-between">
			<h2 class="text-lg font-semibold">Timeline Viewer</h2>
			<div class="text-xs text-slate-400">{{ timelineStore.entries.length }} entries</div>
		</header>
		<div ref="scrollContainer" class="h-[560px] overflow-auto rounded border border-slate-800" @scroll="onScroll">
			<div :style="listStyle">
				<div
					v-for="virtual in visibleItems"
					:key="virtual.index"
					class="absolute left-0 right-0 px-2"
					:style="{ top: `${virtual.index * itemHeight}px` }"
				>
					<TimelineEntry :entry="virtual.item" />
				</div>
			</div>
		</div>
		<p class="text-xs text-slate-400">Rendered from index {{ startIndex }} with virtual scrolling.</p>
	</section>
</template>
