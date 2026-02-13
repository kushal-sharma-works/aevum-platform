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
	<div class="rounded border border-slate-800 p-2 text-xs">
		<div class="flex items-center justify-between">
			<button v-if="isObject" class="text-left text-slate-200" @click="isOpen = !isOpen">
				{{ label ?? 'root' }} {{ isOpen ? '▼' : '▶' }} {{ !isOpen ? summary : '' }}
			</button>
			<span v-else class="text-slate-200">{{ label ?? 'value' }}: {{ String(value) }}</span>
			<button class="text-blue-300" @click="copyValue(value)">Copy</button>
		</div>

		<div v-if="isObject && isOpen" class="mt-2 space-y-2 pl-3">
			<BaseJsonViewer
				v-for="[key, child] in entries"
				:key="key"
				:label="key"
				:value="child"
				:depth="(depth ?? 0) + 1"
			/>
		</div>
	</div>
</template>
