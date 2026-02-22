<template>
	<section class="base-card">
		<slot />
	</section>
</template>

<style scoped>
.base-card {
	border: 1px solid var(--p-content-border-color);
	background: var(--p-content-background);
	border-radius: 0.75rem;
	padding: 1rem;
}
</style>
