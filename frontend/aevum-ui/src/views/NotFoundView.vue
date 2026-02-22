<template>
	<div class="nf-wrap">
		<section class="nf-card">
			<div class="nf-code">404</div>
			<div class="nf-title">Page not found</div>
			<RouterLink to="/dashboard" class="p-button p-component">Go to dashboard</RouterLink>
		</section>
	</div>
</template>

<style scoped>
.nf-wrap {
	min-height: 100vh;
	display: grid;
	place-items: center;
	padding: 1rem;
}

.nf-card {
	width: min(420px, 92vw);
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 1rem;
	display: grid;
	gap: 0.75rem;
	text-align: center;
}

.nf-code {
	font-size: 2rem;
	font-weight: 700;
}

.nf-title {
	font-size: 1rem;
	color: var(--p-text-muted-color);
}
</style>
