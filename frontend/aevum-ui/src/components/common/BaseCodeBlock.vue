<script setup lang="ts">
defineProps<{ code: string }>()
</script>

<template>
	<section class="base-code">
		<pre class="base-pre"><code>{{ code }}</code></pre>
	</section>
</template>

<style scoped>
.base-code {
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.75rem;
}

.base-pre {
	margin: 0;
	overflow: auto;
	font-size: 0.75rem;
}
</style>
