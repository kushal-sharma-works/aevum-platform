<script setup lang="ts">
import type { CausalChain } from '@/types/audit'

defineProps<{ chain: CausalChain }>()
</script>

<template>
	<section class="causal-banner">
		Nodes: {{ chain.nodes.length }} | Edges: {{ chain.edges.length }}
	</section>
</template>

<style scoped>
.causal-banner {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.65rem 0.75rem;
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
}
</style>
