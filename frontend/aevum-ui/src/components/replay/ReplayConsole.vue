<script setup lang="ts">
import { ref } from 'vue'
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
	<q-card flat bordered>
		<q-card-section class="q-gutter-sm">
			<div class="row q-col-gutter-sm">
				<div class="col-12 col-md-4"><q-input v-model="streamId" label="Stream ID" outlined /></div>
				<div class="col-12 col-md-4"><q-input v-model="from" label="From ISO" outlined /></div>
				<div class="col-12 col-md-4"><q-input v-model="to" label="To ISO" outlined /></div>
			</div>
			<div class="row q-gutter-sm">
				<q-btn color="primary" label="Start Replay" @click="start" />
				<q-btn color="negative" label="Stop" @click="replayStore.stopReplay" />
			</div>
		<ReplayProgress :progress="replayStore.replayProgress" />
		<ReplayEventFeed :items="replayStore.replayedEvents" />
		</q-card-section>
	</q-card>
</template>
