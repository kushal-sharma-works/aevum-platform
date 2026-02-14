<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue'

defineProps<{ isReplaying: boolean }>()
const emit = defineEmits<{ (e: 'start'): void; (e: 'stop'): void }>()
</script>

<template>
	<div class="flex gap-2">
		<BaseButton :disabled="isReplaying" @click="emit('start')">Start Replay</BaseButton>
		<BaseButton variant="danger" :disabled="!isReplaying" @click="emit('stop')">Stop Replay</BaseButton>
	</div>
</template>
