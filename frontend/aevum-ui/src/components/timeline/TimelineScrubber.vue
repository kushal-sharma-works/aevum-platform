<script setup lang="ts">
const value = defineModel<number>({ required: true })
</script>

<template>
	<div>
		<input v-model="value" type="range" min="0" max="100" class="w-full" />
	</div>
</template>
