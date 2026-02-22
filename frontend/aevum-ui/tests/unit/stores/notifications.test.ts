import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useNotificationsStore } from '@/stores/notifications'

describe('notifications store', () => {
	beforeEach(() => {
		setActivePinia(createPinia())
		vi.useFakeTimers()
		vi.stubGlobal('crypto', {
			randomUUID: vi.fn(() => 'n1')
		})
	})

	it('adds and auto-removes notification', () => {
		const store = useNotificationsStore()
		store.notify('success', 'saved', 1000)

		expect(store.notifications).toHaveLength(1)
		expect(store.notifications[0]).toMatchObject({ id: 'n1', type: 'success', message: 'saved', duration: 1000 })

		vi.advanceTimersByTime(1000)
		expect(store.notifications).toHaveLength(0)
	})
})
