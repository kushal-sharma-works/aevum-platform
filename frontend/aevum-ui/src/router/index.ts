import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
	history: createWebHistory(),
	routes: [
		{ path: '/', redirect: '/dashboard' },
		{ path: '/dashboard', component: () => import('@/views/DashboardView.vue') },
		{ path: '/timeline', component: () => import('@/views/TimelineView.vue') },
		{ path: '/events', component: () => import('@/views/EventsView.vue') },
		{ path: '/events/:eventId', component: () => import('@/views/EventDetailView.vue') },
		{ path: '/decisions', component: () => import('@/views/DecisionsView.vue') },
		{ path: '/decisions/:decisionId', component: () => import('@/views/DecisionDetailView.vue') },
		{ path: '/rules', component: () => import('@/views/RulesView.vue') },
		{ path: '/rules/create', component: () => import('@/views/RuleCreateView.vue') },
		{ path: '/rules/:ruleId', component: () => import('@/views/RuleDetailView.vue') },
		{ path: '/replay', component: () => import('@/views/ReplayView.vue') },
		{ path: '/diff', component: () => import('@/views/DiffView.vue') },
		{ path: '/audit/:decisionId', component: () => import('@/views/AuditView.vue') },
		{ path: '/login', component: () => import('@/views/LoginView.vue') },
		{ path: '/:pathMatch(.*)*', component: () => import('@/views/NotFoundView.vue') }
	]
})

router.beforeEach((to) => {
	const authStore = useAuthStore()
	if (to.path === '/login' && authStore.isAuthenticated) {
		return '/dashboard'
	}

	if (to.path !== '/login' && !authStore.isAuthenticated) {
		return '/login'
	}

	return true
})

export default router
