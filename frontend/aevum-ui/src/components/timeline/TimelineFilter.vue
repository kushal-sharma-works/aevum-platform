<script setup lang="ts">
const streamId = defineModel<string>('streamId', { required: true })
const type = defineModel<string>('type', { required: true })
</script>

<template>
	<div class="grid gap-2 md:grid-cols-2">
		<input v-model="streamId" class="rounded bg-slate-900 px-2 py-1" placeholder="Stream ID" />
		<input v-model="type" class="rounded bg-slate-900 px-2 py-1" placeholder="Type" />
	</div>
</template>
