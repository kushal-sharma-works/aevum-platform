<script setup lang="ts">
defineProps<{ code: string }>()
</script>

<template>
	<pre class="overflow-auto rounded bg-slate-900 p-3 text-xs text-slate-100"><code>{{ code }}</code></pre>
</template>
