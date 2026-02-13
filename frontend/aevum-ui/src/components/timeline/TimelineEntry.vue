<script setup lang="ts">
import type { TimelineEntry } from '@/types/timeline'

const props = defineProps<{ entry: TimelineEntry }>()

const color = (() => {
	if (props.entry.kind === 'event') {
		return 'bg-blue-600'
	}

	switch (props.entry.item.status) {
		case 'Evaluated':
			return 'bg-emerald-600'
		case 'Skipped':
			return 'bg-amber-600'
		case 'Error':
			return 'bg-rose-600'
	}
})()
</script>

<template>
	<div class="flex items-center gap-2 rounded border border-slate-800 p-2 text-sm">
		<span :class="['inline-block h-2 w-2 rounded-full', color]" />
		<span>{{ entry.timestamp }}</span>
		<span class="text-slate-400">{{ entry.kind }}</span>
	</div>
</template>
