<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppFooter from '@/components/layout/AppFooter.vue'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'

const route = useRoute()
const showShell = computed(() => route.path !== '/login')
</script>

<template>
	<div class="dark min-h-screen bg-slate-950 text-slate-100">
		<template v-if="showShell">
			<AppHeader />
			<div class="flex">
				<AppSidebar />
				<main class="min-h-[calc(100vh-8rem)] flex-1 p-4">
					<RouterView />
				</main>
			</div>
			<AppFooter />
		</template>
		<RouterView v-else />
	</div>
</template>
