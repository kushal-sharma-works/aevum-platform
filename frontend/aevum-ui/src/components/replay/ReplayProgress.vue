<script setup lang="ts">
defineProps<{ progress: number }>()
</script>

<template>
	<div class="replay-progress">
		<div class="progress-caption">Progress {{ progress.toFixed(0) }}%</div>
		<div class="progress-track">
			<div class="progress-fill" :style="{ width: `${progress}%` }" />
		</div>
	</div>
</template>

<style scoped>
.replay-progress {
	display: grid;
	gap: 0.35rem;
}

.progress-caption {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}

.progress-track {
	height: 8px;
	border-radius: 999px;
	background: var(--p-surface-300);
	overflow: hidden;
}

.progress-fill {
	height: 100%;
	background: var(--p-primary-color);
	transition: width 0.2s ease;
}
</style>
