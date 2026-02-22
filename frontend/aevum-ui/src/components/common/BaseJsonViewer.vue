<script setup lang="ts">
import { computed, ref } from 'vue'

const props = defineProps<{ value: unknown; label?: string; depth?: number }>()
const isOpen = ref((props.depth ?? 0) < 1)

const isObject = computed(() => typeof props.value === 'object' && props.value !== null)

const entries = computed(() => {
	if (!isObject.value) {
		return [] as Array<[string, unknown]>
	}

	if (Array.isArray(props.value)) {
		return props.value.map((item, index) => [String(index), item] as [string, unknown])
	}

	return Object.entries(props.value as Record<string, unknown>)
})

const summary = computed(() => {
	if (!isObject.value) {
		return ''
	}
	return Array.isArray(props.value) ? `Array(${entries.value.length})` : `Object(${entries.value.length})`
})

function copyValue(value: unknown): void {
	navigator.clipboard.writeText(JSON.stringify(value, null, 2)).catch(() => undefined)
}
</script>

<template>
	<section class="json-card">
		<div class="json-header">
			<button
				v-if="isObject"
				type="button"
				class="p-button p-button-text"
				@click="isOpen = !isOpen"
			>
				{{ label ?? 'root' }} {{ isOpen ? '▼' : '▶' }} {{ !isOpen ? summary : '' }}
			</button>
			<span v-else>{{ label ?? 'value' }}: {{ String(value) }}</span>
			<button type="button" class="p-button p-button-text" @click="copyValue(value)">Copy</button>
		</div>

		<div v-if="isObject && isOpen" class="json-children">
			<BaseJsonViewer
				v-for="[key, child] in entries"
				:key="key"
				:label="key"
				:value="child"
				:depth="(depth ?? 0) + 1"
			/>
		</div>
	</section>
</template>

<style scoped>
.json-card {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	padding: 0.5rem 0.75rem;
	font-size: 0.8rem;
}

.json-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 0.5rem;
}

.json-children {
	margin-top: 0.5rem;
	padding-left: 0.75rem;
	display: grid;
	gap: 0.5rem;
}
</style>
