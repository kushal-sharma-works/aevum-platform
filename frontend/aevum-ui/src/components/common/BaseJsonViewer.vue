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
	<q-card flat bordered class="text-caption">
		<q-card-section class="q-py-sm q-px-md">
			<div class="row items-center justify-between">
				<q-btn
					v-if="isObject"
					flat
					dense
					color="grey-8"
					:label="`${label ?? 'root'} ${isOpen ? '▼' : '▶'} ${!isOpen ? summary : ''}`"
					@click="isOpen = !isOpen"
				/>
				<span v-else>{{ label ?? 'value' }}: {{ String(value) }}</span>
				<q-btn flat dense color="primary" label="Copy" @click="copyValue(value)" />
			</div>

			<div v-if="isObject && isOpen" class="q-mt-sm q-gutter-sm" style="padding-left: 12px">
			<BaseJsonViewer
				v-for="[key, child] in entries"
				:key="key"
				:label="key"
				:value="child"
				:depth="(depth ?? 0) + 1"
			/>
		</div>
		</q-card-section>
	</q-card>
</template>
