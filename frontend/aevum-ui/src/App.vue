<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const showShell = computed(() => route.path !== '/login')

const navLinks = [
	{ label: 'Dashboard', to: '/dashboard', icon: 'dashboard' },
	{ label: 'Timeline', to: '/timeline', icon: 'timeline' },
	{ label: 'Events', to: '/events', icon: 'event' },
	{ label: 'Decisions', to: '/decisions', icon: 'account_tree' },
	{ label: 'Rules', to: '/rules', icon: 'rule' },
	{ label: 'Replay', to: '/replay', icon: 'play_circle' },
	{ label: 'Diff', to: '/diff', icon: 'difference' }
]
</script>

<template>
	<q-layout view="lHh Lpr lFf">
		<template v-if="showShell">
			<q-header elevated>
				<q-toolbar>
					<q-toolbar-title>Aevum Platform</q-toolbar-title>
					<q-badge color="secondary" rounded>SIT</q-badge>
				</q-toolbar>
			</q-header>

			<q-drawer show-if-above bordered :width="240">
				<q-list>
					<q-item-label header>Navigation</q-item-label>
					<q-item
						v-for="link in navLinks"
						:key="link.to"
						:to="link.to"
						clickable
					>
						<q-item-section avatar>
							<q-icon :name="link.icon" />
						</q-item-section>
						<q-item-section>{{ link.label }}</q-item-section>
					</q-item>
				</q-list>
			</q-drawer>

			<q-page-container>
				<q-page class="q-pa-md bg-grey-2">
					<RouterView />
				</q-page>
			</q-page-container>

			<q-footer bordered class="bg-white text-grey-8">
				<q-toolbar>
					<q-toolbar-title class="text-subtitle2">Aevum UI · Quasar</q-toolbar-title>
				</q-toolbar>
			</q-footer>
		</template>

		<q-page-container v-else>
			<q-page class="q-pa-md bg-grey-2">
				<RouterView />
			</q-page>
		</q-page-container>
	</q-layout>
</template>
