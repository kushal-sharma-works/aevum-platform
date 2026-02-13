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
	<aside class="hidden w-56 border-r border-slate-800 p-3 md:block">
		<ul class="space-y-2 text-sm">
			<li v-for="link in links" :key="link.to">
				<RouterLink :to="link.to" class="block rounded px-3 py-2 text-slate-300 hover:bg-slate-900 hover:text-white">
					{{ link.label }}
				</RouterLink>
			</li>
		</ul>
	</aside>
</template>
