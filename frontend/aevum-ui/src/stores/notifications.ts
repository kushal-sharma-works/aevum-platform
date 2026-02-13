import { ref } from 'vue'
import { defineStore } from 'pinia'

export type NotificationType = 'success' | 'error' | 'warning' | 'info'

export interface NotificationItem {
	readonly id: string
	readonly type: NotificationType
	readonly message: string
	readonly duration: number
}

export const useNotificationsStore = defineStore('notifications', () => {
	const notifications = ref<NotificationItem[]>([])

	function remove(id: string): void {
		notifications.value = notifications.value.filter((item) => item.id !== id)
	}

	function notify(type: NotificationType, message: string, duration = 5000): void {
		const id = crypto.randomUUID()
		notifications.value.push({ id, type, message, duration })
		window.setTimeout(() => {
			remove(id)
		}, duration)
	}

	return {
		notifications,
		notify,
		remove
	}
})
