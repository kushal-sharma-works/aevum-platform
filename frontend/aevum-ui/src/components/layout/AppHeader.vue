<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
</script>

<template>
	<header class="flex h-16 items-center justify-between border-b border-slate-800 px-4">
		<div class="text-lg font-semibold">Aevum Platform</div>
		<nav class="flex items-center gap-4 text-sm text-slate-300">
			<RouterLink to="/dashboard" class="hover:text-white">Dashboard</RouterLink>
			<RouterLink to="/timeline" class="hover:text-white">Timeline</RouterLink>
			<RouterLink to="/decisions" class="hover:text-white">Decisions</RouterLink>
			<button class="rounded bg-slate-800 px-3 py-1 hover:bg-slate-700" @click="authStore.logout">Logout</button>
		</nav>
	</header>
</template>
