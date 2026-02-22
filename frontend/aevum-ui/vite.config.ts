import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
	plugins: [vue()],
	resolve: {
		alias: {
			'@': fileURLToPath(new URL('./src', import.meta.url))
		}
	},
	server: {
		port: 3000,
		proxy: {
			'/api/events/admin': {
				target: 'http://localhost:9091',
				changeOrigin: true,
				rewrite: (path) => path.replace('/api/events/admin', '/admin')
			},
			'/api/events': {
				target: 'http://localhost:8081',
				changeOrigin: true,
				rewrite: (path) => path.replace('/api/events', '/api/v1')
			},
			'/api/decisions': {
				target: 'http://localhost:8080',
				changeOrigin: true,
				rewrite: (path) => path.replace('/api/decisions', '/api/v1')
			},
			'/api/query': {
				target: 'http://localhost:8082',
				changeOrigin: true,
				rewrite: (path) => path.replace('/api/query', '/api/v1')
			}
		}
	},
	test: {
		environment: 'jsdom',
		globals: true,
		include: ['tests/**/*.test.ts'],
		setupFiles: ['tests/setup.ts']
	}
})
