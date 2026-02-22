<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ message: string; type?: 'success' | 'error' | 'warning' | 'info' }>()

const toneClass = computed(() => {
	if (props.type === 'success') return 'p-message-success'
	if (props.type === 'error') return 'p-message-error'
	if (props.type === 'warning') return 'p-message-warn'
	return 'p-message-info'
})
</script>

<template>
	<div class="p-message p-component" :class="toneClass">
		{{ message }}
	</div>
</template>
