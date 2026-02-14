<script setup lang="ts">
defineProps<{ step: { readonly timestamp: string; readonly message: string } }>()
</script>

<template>
	<q-card flat bordered>
		<q-card-section class="q-py-sm q-px-md">
			<div class="text-caption text-grey-7">{{ step.timestamp }}</div>
			<div class="text-body2">{{ step.message }}</div>
		</q-card-section>
	</q-card>
</template>
