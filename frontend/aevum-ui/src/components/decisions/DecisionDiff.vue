<script setup lang="ts">
import { computed } from 'vue'
import BaseJsonViewer from '@/components/common/BaseJsonViewer.vue'
import { diffJson } from '@/utils/jsonDiff'

const props = defineProps<{ left: Record<string, unknown>; right: Record<string, unknown> }>()
const diffs = computed(() => diffJson(props.left, props.right))
</script>

<template>
	<section class="grid gap-4 lg:grid-cols-2">
		<div class="space-y-2 rounded border border-slate-800 p-3">
			<h3 class="font-semibold">T1 / V1</h3>
			<BaseJsonViewer :value="left" />
		</div>
		<div class="space-y-2 rounded border border-slate-800 p-3">
			<h3 class="font-semibold">T2 / V2</h3>
			<BaseJsonViewer :value="right" />
		</div>
		<div class="lg:col-span-2">
			<h4 class="mb-2 font-semibold">Changes</h4>
			<ul class="space-y-1 text-xs">
				<li
					v-for="entry in diffs"
					:key="entry.path"
					:class="[
						'rounded p-2',
						entry.type === 'added'
							? 'bg-emerald-900/30'
							: entry.type === 'removed'
								? 'bg-rose-900/30'
								: 'bg-amber-900/30'
					]"
				>
					{{ entry.type }} @ {{ entry.path }}
				</li>
			</ul>
		</div>
	</section>
</template>
