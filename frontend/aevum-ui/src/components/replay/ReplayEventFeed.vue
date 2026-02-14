<script setup lang="ts">
import type { Decision } from '@/types/decision'
import type { Event } from '@/types/event'

defineProps<{ items: ReadonlyArray<Event | Decision> }>()
</script>

<template>
	<q-list bordered separator class="feed-list">
		<q-item v-for="(item, index) in items" :key="index">
			<q-item-section>
				<q-item-label>{{ 'decisionId' in item ? `Decision ${item.decisionId}` : `Event ${item.eventId}` }}</q-item-label>
			</q-item-section>
		</q-item>
	</q-list>
</template>

<style scoped>
.feed-list {
	max-height: 240px;
	overflow: auto;
}
</style>
