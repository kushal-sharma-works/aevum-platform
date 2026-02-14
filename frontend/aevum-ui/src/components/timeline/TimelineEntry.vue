<script setup lang="ts">
import { computed } from 'vue'
import type { TimelineEntry } from '@/types/timeline'

const props = defineProps<{ entry: TimelineEntry }>()

const color = computed(() => {
	if (props.entry.kind === 'event') {
		return 'primary'
	}

	switch (props.entry.item.status) {
		case 'Evaluated':
			return 'positive'
		case 'Skipped':
			return 'warning'
		default:
			return 'negative'
	}
})
</script>

<template>
	<q-item dense class="rounded-borders bg-white q-px-sm q-py-xs">
		<q-item-section side>
			<q-badge rounded :color="color" />
		</q-item-section>
		<q-item-section>
			<q-item-label>{{ entry.timestamp }}</q-item-label>
			<q-item-label caption>{{ entry.kind }}</q-item-label>
		</q-item-section>
	</q-item>
</template>
