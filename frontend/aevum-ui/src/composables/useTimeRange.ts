import { computed, ref } from 'vue'

export type TimePreset = 'Last 1h' | 'Last 24h' | 'Last 7d' | 'Last 30d' | 'Custom'

export function useTimeRange() {
	const to = ref(new Date())
	const from = ref(new Date(Date.now() - 60 * 60 * 1000))
	const preset = ref<TimePreset>('Last 1h')

	const isValid = computed(() => from.value.getTime() < to.value.getTime())

	function setPreset(value: TimePreset): void {
		preset.value = value
		const now = new Date()
		to.value = now
		switch (value) {
			case 'Last 1h':
				from.value = new Date(now.getTime() - 60 * 60 * 1000)
				break
			case 'Last 24h':
				from.value = new Date(now.getTime() - 24 * 60 * 60 * 1000)
				break
			case 'Last 7d':
				from.value = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
				break
			case 'Last 30d':
				from.value = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000)
				break
			case 'Custom':
				break
		}
	}

	function setCustomRange(nextFrom: Date, nextTo: Date): void {
		if (nextFrom.getTime() >= nextTo.getTime()) {
			return
		}
		preset.value = 'Custom'
		from.value = nextFrom
		to.value = nextTo
	}

	return { from, to, preset, isValid, setPreset, setCustomRange }
}
