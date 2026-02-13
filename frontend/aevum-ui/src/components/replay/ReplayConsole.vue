<script setup lang="ts">
import { ref } from 'vue'
import BaseButton from '@/components/common/BaseButton.vue'
import { useReplayStore } from '@/stores/replay'
import ReplayEventFeed from './ReplayEventFeed.vue'
import ReplayProgress from './ReplayProgress.vue'

const replayStore = useReplayStore()
const streamId = ref('')
const from = ref(new Date(Date.now() - 3600_000).toISOString())
const to = ref(new Date().toISOString())

const start = async () => {
	await replayStore.startReplay(streamId.value, from.value, to.value)
}
</script>

<template>
	<div class="space-y-3 rounded border border-slate-800 p-3">
		<div class="grid gap-2 md:grid-cols-3">
			<input v-model="streamId" class="rounded bg-slate-900 px-2 py-1" placeholder="Stream ID" />
			<input v-model="from" class="rounded bg-slate-900 px-2 py-1" placeholder="From ISO" />
			<input v-model="to" class="rounded bg-slate-900 px-2 py-1" placeholder="To ISO" />
		</div>
		<div class="flex gap-2">
			<BaseButton @click="start">Start Replay</BaseButton>
			<BaseButton variant="danger" @click="replayStore.stopReplay">Stop</BaseButton>
		</div>
		<ReplayProgress :progress="replayStore.replayProgress" />
		<ReplayEventFeed :items="replayStore.replayedEvents" />
	</div>
</template>
