<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ label: string; tone?: 'success' | 'warning' | 'error' | 'info' }>()

const color = computed(() => {
	if (props.tone === 'success') return 'p-tag-success'
	if (props.tone === 'warning') return 'p-tag-warn'
	if (props.tone === 'error') return 'p-tag-danger'
	return 'p-tag-info'
})
</script>

<template>
	<span class="p-tag p-component" :class="color">{{ label }}</span>
</template>
