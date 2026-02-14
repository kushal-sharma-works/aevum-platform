<script setup lang="ts">
defineProps<{ title: string; description?: string }>()
</script>

<template>
	<section class="q-gutter-md">
		<header>
			<div class="text-h5 text-weight-medium">{{ title }}</div>
			<div v-if="description" class="text-subtitle2 text-grey-7">{{ description }}</div>
		</header>
		<q-separator />
		<div>
			<slot />
		</div>
	</section>
</template>
