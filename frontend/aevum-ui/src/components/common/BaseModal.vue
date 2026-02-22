<script setup lang="ts">
defineProps<{ open: boolean; title: string }>()
const emit = defineEmits<{ (e: 'close'): void }>()
</script>

<template>
	<div v-if="open" class="modal-backdrop" @click.self="emit('close')">
		<div class="modal-card">
			<header class="modal-header">
				<div class="modal-title">{{ title }}</div>
				<button type="button" class="p-button p-button-text p-button-rounded" aria-label="Close" @click="emit('close')">
					<i class="pi pi-times" />
				</button>
			</header>
			<div class="modal-content">
				<slot />
			</div>
		</div>
	</div>
</template>

<style scoped>
.modal-backdrop {
	position: fixed;
	inset: 0;
	background: rgba(2, 6, 23, 0.45);
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 1rem;
	z-index: 1100;
}

.modal-card {
	width: min(90vw, 680px);
	background: var(--p-content-background);
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	overflow: hidden;
}

.modal-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 0.75rem 1rem;
	border-bottom: 1px solid var(--p-content-border-color);
}

.modal-title {
	font-size: 1.1rem;
	font-weight: 600;
}

.modal-content {
	padding: 1rem;
}
</style>
