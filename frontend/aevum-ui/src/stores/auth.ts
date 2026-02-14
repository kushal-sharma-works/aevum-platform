import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

interface Credentials {
	readonly username: string
	readonly password: string
}

export const useAuthStore = defineStore('auth', () => {
	const token = ref<string | null>(localStorage.getItem('aevum_token'))
	const isLoading = ref(false)
	const isAuthenticated = computed(() => Boolean(token.value))

	async function login(credentials: Credentials): Promise<void> {
		isLoading.value = true
		try {
			const nextToken = btoa(`${credentials.username}:${credentials.password}:${Date.now()}`)
			token.value = nextToken
			localStorage.setItem('aevum_token', nextToken)
		} finally {
			isLoading.value = false
		}
	}

	function logout(): void {
		token.value = null
		localStorage.removeItem('aevum_token')
		if (window.location.pathname !== '/login') {
			window.location.href = '/login'
		}
	}

	return {
		token,
		isLoading,
		isAuthenticated,
		login,
		logout
	}
})
