import { createPinia, setActivePinia } from 'pinia'
import { describe, expect, it } from 'vitest'
import { useAuthStore } from '@/stores/auth'

describe('auth store', () => {
	it('sets token on login', async () => {
		setActivePinia(createPinia())
		const store = useAuthStore()
		await store.login({ username: 'u', password: 'p' })
		expect(store.isAuthenticated).toBe(true)
	})
})
