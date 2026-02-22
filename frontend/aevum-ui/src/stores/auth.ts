import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { resetUnauthorizedHandling } from '@/api/client'

function isJwtLikeToken(value: string): boolean {
	const parts = value.split('.')
	return parts.length === 3 && parts.every((part) => part.length > 0)
}

function toBase64Url(input: string): string {
	return btoa(input).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
}

function bytesToBase64Url(bytes: Uint8Array): string {
	let binary = ''
	for (const byte of bytes) {
		binary += String.fromCharCode(byte)
	}
	return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
}

async function createDevJwt(username: string): Promise<string | null> {
	if (!window.isSecureContext || !window.crypto?.subtle) {
		return null
	}

	const secret = import.meta.env.VITE_AUTH_JWT_SECRET ?? 'dev-local-secret'
	const now = Math.floor(Date.now() / 1000)
	const header = toBase64Url(JSON.stringify({ alg: 'HS256', typ: 'JWT' }))
	const payload = toBase64Url(
		JSON.stringify({
			iss: 'aevum-ui-dev',
			sub: username,
			role: 'developer',
			iat: now,
			exp: now + 8 * 60 * 60
		})
	)

	const content = `${header}.${payload}`
	const key = await window.crypto.subtle.importKey('raw', new TextEncoder().encode(secret), { name: 'HMAC', hash: 'SHA-256' }, false, [
		'sign'
	])
	const signature = await window.crypto.subtle.sign('HMAC', key, new TextEncoder().encode(content))
	return `${content}.${bytesToBase64Url(new Uint8Array(signature))}`
}

interface Credentials {
	readonly username: string
	readonly password: string
}

export const useAuthStore = defineStore('auth', () => {
	const initialToken = localStorage.getItem('aevum_token')
	const token = ref<string | null>(initialToken && isJwtLikeToken(initialToken) ? initialToken : null)
	const isLoading = ref(false)
	const isAuthenticated = computed(() => Boolean(token.value && isJwtLikeToken(token.value)))

	if (initialToken && !isJwtLikeToken(initialToken)) {
		localStorage.removeItem('aevum_token')
	}

	async function login(credentials: Credentials): Promise<void> {
		isLoading.value = true
		try {
			const trimmedUsername = credentials.username.trim()
			const trimmedPassword = credentials.password.trim()
			if (!trimmedUsername || !trimmedPassword) {
				throw new Error('Username and password are required')
			}

			const nextToken = await createDevJwt(trimmedUsername)
			if (!nextToken) {
				throw new Error('Unable to create local JWT for this browser context')
			}
			token.value = nextToken
			localStorage.setItem('aevum_token', nextToken)
			resetUnauthorizedHandling()
		} finally {
			isLoading.value = false
		}
	}

	function logout(): void {
		token.value = null
		localStorage.removeItem('aevum_token')
		resetUnauthorizedHandling()
		if (window.location.pathname !== '/login') {
			window.dispatchEvent(new CustomEvent('aevum:auth-required'))
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
