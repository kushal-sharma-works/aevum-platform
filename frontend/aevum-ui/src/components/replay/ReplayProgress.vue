<script setup lang="ts">
defineProps<{ progress: number }>()
</script>

<template>
	<div class="q-gutter-xs">
		<div class="text-caption text-grey-7">Progress {{ progress.toFixed(0) }}%</div>
		<q-linear-progress :value="progress / 100" rounded size="8px" color="primary" track-color="grey-4" />
	</div>
</template>
