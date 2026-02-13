<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseButton from '@/components/common/BaseButton.vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const username = ref('')
const password = ref('')

const submit = async () => {
	await authStore.login({ username: username.value, password: password.value })
	await router.push('/dashboard')
}
</script>

<template>
	<main class="flex min-h-screen items-center justify-center">
		<form class="w-full max-w-sm space-y-3 rounded border border-slate-800 p-5" @submit.prevent="submit">
			<h1 class="text-lg font-semibold">Aevum Login</h1>
			<input v-model="username" class="w-full rounded bg-slate-900 px-3 py-2" placeholder="Username" />
			<input v-model="password" type="password" class="w-full rounded bg-slate-900 px-3 py-2" placeholder="Password" />
			<BaseButton type="submit" class="w-full">Login</BaseButton>
		</form>
	</main>
</template>
