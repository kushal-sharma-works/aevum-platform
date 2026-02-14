<script setup lang="ts">
import type { CausalChain } from '@/types/audit'

defineProps<{ chain: CausalChain }>()
</script>

<template>
	<div class="rounded border border-slate-800 p-3 text-xs">
		Nodes: {{ chain.nodes.length }} | Edges: {{ chain.edges.length }}
	</div>
</template>
