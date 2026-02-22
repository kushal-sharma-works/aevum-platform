import { createPinia, setActivePinia } from 'pinia'
import { mount } from '@vue/test-utils'
import { describe, expect, it, vi } from 'vitest'
import AppFooter from '@/components/layout/AppFooter.vue'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import PageContainer from '@/components/layout/PageContainer.vue'
import { useAuthStore } from '@/stores/auth'

const routerLinkStub = {
	props: ['to'],
	template: '<a :data-to="to"><slot /></a>'
}

describe('layout components', () => {
	it('renders footer and page container content', () => {
		const footer = mount(AppFooter)
		expect(footer.text()).toContain('Deterministic Decision Platform')

		const page = mount(PageContainer, {
			props: { title: 'Timeline', description: 'Recent activity' },
			slots: { default: '<div>slot content</div>' },
			global: { stubs: { 'q-separator': true } }
		})
		expect(page.text()).toContain('Timeline')
		expect(page.text()).toContain('slot content')
	})

	it('renders navigation links and triggers logout', async () => {
		setActivePinia(createPinia())
		const authStore = useAuthStore()
		const logoutSpy = vi.spyOn(authStore, 'logout').mockImplementation(() => undefined)

		const header = mount(AppHeader, {
			global: { stubs: { RouterLink: routerLinkStub } }
		})
		expect(header.text()).toContain('Aevum Platform')
		const logoutButton = header
			.findAll('button')
			.find((button) => button.text().includes('Logout'))
		expect(logoutButton).toBeDefined()
		await logoutButton!.trigger('click')
		expect(logoutSpy).toHaveBeenCalledTimes(1)

		const sidebar = mount(AppSidebar, {
			global: { stubs: { RouterLink: routerLinkStub } }
		})
		expect(sidebar.text()).toContain('Dashboard')
		expect(sidebar.text()).toContain('Diff')
	})
})
