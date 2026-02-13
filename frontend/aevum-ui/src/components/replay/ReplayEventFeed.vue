<script setup lang="ts">
import type { Decision } from '@/types/decision'
import type { Event } from '@/types/event'

defineProps<{ items: ReadonlyArray<Event | Decision> }>()
</script>

<template>
	<div class="max-h-60 space-y-1 overflow-auto rounded border border-slate-800 p-2 text-xs">
		<div v-for="(item, index) in items" :key="index" class="rounded bg-slate-900 px-2 py-1">
			{{ 'decisionId' in item ? `Decision ${item.decisionId}` : `Event ${item.eventId}` }}
		</div>
	</div>
</template>
