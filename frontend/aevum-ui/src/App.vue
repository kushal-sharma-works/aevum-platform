<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const showShell = computed(() => route.path !== '/login')
const environmentLabel = (import.meta.env.VITE_APP_ENV ?? 'Local').toUpperCase()

const navLinks = [
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
	<div class="app-shell" v-if="showShell">
		<header class="app-header">
			<div class="brand">Aevum Platform</div>
			<div class="header-actions">
				<span class="p-tag p-component">{{ environmentLabel }}</span>
				<button type="button" class="p-button p-component" @click="authStore.logout">Logout</button>
			</div>
		</header>

		<div class="app-content-wrap">
			<aside class="app-sidebar">
				<div class="sidebar-title">Navigation</div>
				<nav class="sidebar-nav">
					<RouterLink v-for="link in navLinks" :key="link.to" :to="link.to" class="p-button p-button-text nav-link">
						{{ link.label }}
					</RouterLink>
				</nav>
			</aside>

			<main class="app-main">
				<RouterView />
			</main>
		</div>

		<footer class="app-footer">Aevum UI · PrimeVue</footer>
	</div>

	<main class="app-main" v-else>
		<RouterView />
	</main>
</template>

<style scoped>
.app-shell {
	min-height: 100vh;
	display: flex;
	flex-direction: column;
}

.app-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 0.75rem 1rem;
	background: var(--p-content-background);
	border-bottom: 1px solid var(--p-content-border-color);
}

.header-actions {
	display: flex;
	align-items: center;
	gap: 0.5rem;
}

.brand {
	font-size: 1.05rem;
	font-weight: 600;
}

.app-content-wrap {
	display: flex;
	flex: 1;
}

.app-sidebar {
	width: 240px;
	background: var(--p-content-background);
	border-right: 1px solid var(--p-content-border-color);
	padding: 1rem 0.75rem;
}

.sidebar-title {
	font-size: 0.8rem;
	color: var(--p-text-muted-color);
	padding: 0 0.5rem 0.5rem;
}

.sidebar-nav {
	display: flex;
	flex-direction: column;
	gap: 0.25rem;
}

.nav-link {
	justify-content: flex-start;
	text-decoration: none;
	border-radius: 0.5rem;
}

.nav-link.router-link-active {
	background: var(--p-primary-100);
	color: var(--p-primary-700);
}

.app-main {
	flex: 1;
	padding: 1rem;
}

.app-footer {
	padding: 0.75rem 1rem;
	font-size: 0.85rem;
	color: var(--p-text-muted-color);
	background: var(--p-content-background);
	border-top: 1px solid var(--p-content-border-color);
}
</style>
