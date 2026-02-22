<script setup lang="ts">
import { ref } from 'vue'
import { useReplayStore } from '@/stores/replay'
import ReplayEventFeed from './ReplayEventFeed.vue'
import ReplayProgress from './ReplayProgress.vue'

const replayStore = useReplayStore()
const streamId = ref('default')
const from = ref(new Date(Date.now() - 3600_000).toISOString())
const to = ref(new Date().toISOString())

const start = async () => {
	await replayStore.startReplay(streamId.value, from.value, to.value)
}
</script>

<template>
	<section class="replay-card">
		<div v-if="replayStore.error" class="replay-error">{{ replayStore.error }}</div>
		<div class="replay-grid">
			<label class="replay-field">
				<span>Stream ID</span>
				<input v-model="streamId" class="p-inputtext p-component" />
			</label>
			<label class="replay-field">
				<span>From ISO</span>
				<input v-model="from" class="p-inputtext p-component" />
			</label>
			<label class="replay-field">
				<span>To ISO</span>
				<input v-model="to" class="p-inputtext p-component" />
			</label>
		</div>
		<div class="replay-actions">
			<button type="button" class="p-button p-component" :disabled="replayStore.isLoading" @click="start">Start Replay</button>
			<button type="button" class="p-button p-component p-button-danger" :disabled="!replayStore.isReplaying" @click="replayStore.stopReplay">Stop</button>
		</div>
		<ReplayProgress :progress="replayStore.replayProgress" />
		<ReplayEventFeed :items="replayStore.replayedEvents" />
	</section>
</template>

<style scoped>
.replay-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 1rem;
	display: grid;
	gap: 0.75rem;
}

.replay-grid {
	display: grid;
	grid-template-columns: repeat(3, minmax(0, 1fr));
	gap: 0.75rem;
}

@media (max-width: 900px) {
	.replay-grid {
		grid-template-columns: 1fr;
	}
}

.replay-field {
	display: grid;
	gap: 0.35rem;
}

.replay-field span {
	font-size: 0.85rem;
	color: var(--p-text-muted-color);
}

.replay-actions {
	display: flex;
	gap: 0.5rem;
}

.replay-error {
	padding: 0.55rem 0.7rem;
	border-radius: 0.5rem;
	border: 1px solid color-mix(in oklab, var(--p-red-500) 40%, var(--p-content-border-color));
	background: color-mix(in oklab, var(--p-red-500) 10%, transparent);
	color: var(--p-red-500);
	font-size: 0.85rem;
}
</style>
