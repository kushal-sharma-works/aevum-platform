import axios, { AxiosError, type InternalAxiosRequestConfig } from 'axios'
import type { ApiError } from '@/types/common'

const RETRY_KEY = 'x-aevum-retry-count'

const client = axios.create({
	timeout: 15000,
	headers: {
		'Content-Type': 'application/json'
	}
})

client.interceptors.request.use((config: InternalAxiosRequestConfig) => {
	const token = localStorage.getItem('aevum_token')
	if (token) {
		config.headers.Authorization = `Bearer ${token}`
	}
	return config
})

client.interceptors.response.use(
	(response) => response,
	async (error: AxiosError) => {
		const status = error.response?.status
		const config = error.config

		if (status === 401) {
			localStorage.removeItem('aevum_token')
			if (window.location.pathname !== '/login') {
				window.location.href = '/login'
			}
		}

		if (status !== undefined && status >= 500 && config) {
			const retries = Number(config.headers?.[RETRY_KEY] ?? 0)
			if (retries < 2) {
				await new Promise((resolve) => {
					window.setTimeout(resolve, 1000)
				})
				config.headers[RETRY_KEY] = String(retries + 1)
				return client(config)
			}
		}

		const payload = error.response?.data as Partial<ApiError> | undefined
		const normalized: ApiError = {
			type: payload?.type ?? 'about:blank',
			title: payload?.title ?? 'Request failed',
			status: payload?.status ?? status ?? 500,
			detail: payload?.detail ?? error.message,
			traceId: payload?.traceId
		}

		return Promise.reject(normalized)
	}
)

export default client
