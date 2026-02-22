<script setup lang="ts">
const links: ReadonlyArray<{ readonly label: string; readonly to: string }> = [
	{ label: 'Dashboard', to: '/dashboard' },
	{ label: 'Timeline', to: '/timeline' },
	{ label: 'Events', to: '/events' },
	{ label: 'Decisions', to: '/decisions' },
	{ label: 'Rules', to: '/rules' },
	{ label: 'Replay', to: '/replay' },
	{ label: 'Diff', to: '/diff' }
]
</script>

<template>
	<nav class="sidebar-nav">
		<RouterLink v-for="link in links" :key="link.to" :to="link.to" class="sidebar-link">{{ link.label }}</RouterLink>
	</nav>
</template>

<style scoped>
.sidebar-nav {
	display: grid;
	gap: 0.3rem;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 0.5rem;
}

.sidebar-link {
	padding: 0.45rem 0.55rem;
	border-radius: 0.5rem;
	text-decoration: none;
	color: var(--p-text-color);
	font-weight: 500;
}

.sidebar-link.router-link-active {
	background: color-mix(in oklab, var(--p-primary-color) 14%, transparent);
	color: var(--p-primary-color);
}
</style>
