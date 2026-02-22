<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<template>
	<header class="app-header">
		<div class="app-title">Aevum Platform</div>
		<nav class="app-nav">
			<RouterLink to="/dashboard" class="app-link">Dashboard</RouterLink>
			<RouterLink to="/timeline" class="app-link">Timeline</RouterLink>
			<RouterLink to="/decisions" class="app-link">Decisions</RouterLink>
		</nav>
		<button type="button" class="p-button p-component" @click="authStore.logout">Logout</button>
	</header>
</template>

<style scoped>
.app-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 1rem;
	padding: 0.75rem 1rem;
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
}

.app-title {
	font-size: 1.1rem;
	font-weight: 700;
}

.app-nav {
	display: flex;
	gap: 0.4rem;
	margin-left: auto;
	margin-right: 0.5rem;
}

.app-link {
	text-decoration: none;
	padding: 0.4rem 0.65rem;
	border-radius: 0.5rem;
	color: var(--p-text-color);
	font-weight: 500;
}

.app-link.router-link-active {
	background: color-mix(in oklab, var(--p-primary-color) 14%, transparent);
	color: var(--p-primary-color);
}
</style>
