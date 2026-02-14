<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'

const props = defineProps<{ open: boolean; title: string }>()
const emit = defineEmits<{ (e: 'close'): void }>()

const onKeyDown = (event: KeyboardEvent) => {
	if (props.open && event.key === 'Escape') {
		emit('close')
	}
}

onMounted(() => {
	window.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
	window.removeEventListener('keydown', onKeyDown)
})
</script>

<template>
	<Teleport to="body">
		<div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60" @click.self="emit('close')">
			<div class="w-full max-w-lg rounded bg-slate-900 p-4 shadow-lg">
				<div class="mb-3 flex items-center justify-between">
					<h2 class="text-lg font-semibold">{{ title }}</h2>
					<button class="text-slate-300" @click="emit('close')">×</button>
				</div>
				<slot />
			</div>
		</div>
	</Teleport>
</template>
