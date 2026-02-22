<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const username = ref('')
const password = ref('')
const loginError = ref<string | null>(null)

const submit = async () => {
	loginError.value = null
	try {
		await authStore.login({ username: username.value, password: password.value })
		await router.push('/dashboard')
	} catch (error) {
		loginError.value = error instanceof Error ? error.message : 'Login failed'
	}
}
</script>

<template>
	<div class="login-wrap">
		<section class="login-card">
			<header>
				<div class="login-title">Aevum Login</div>
				<div class="login-subtitle">Sign in to continue</div>
			</header>
			<form class="login-form" @submit.prevent="submit">
				<label class="login-field">
					<span>Username</span>
					<input v-model="username" autocomplete="username" class="p-inputtext p-component" />
				</label>
				<label class="login-field">
					<span>Password</span>
					<input v-model="password" type="password" autocomplete="current-password" class="p-inputtext p-component" />
				</label>
				<button type="submit" class="p-button p-component" :disabled="authStore.isLoading">Sign In</button>
				<div v-if="loginError" class="login-error">{{ loginError }}</div>
			</form>
		</section>
	</div>
</template>

<style scoped>
.login-wrap {
	min-height: 100vh;
	display: grid;
	place-items: center;
	padding: 1rem;
}

.login-card {
	width: min(460px, 92vw);
	border: 1px solid var(--p-content-border-color);
	border-radius: 0.75rem;
	background: var(--p-content-background);
	padding: 1rem;
	display: grid;
	gap: 1rem;
}

.login-title {
	font-size: 1.15rem;
	font-weight: 600;
}

.login-subtitle {
	font-size: 0.9rem;
	color: var(--p-text-muted-color);
}

.login-form {
	display: grid;
	gap: 0.75rem;
}

.login-field {
	display: grid;
	gap: 0.35rem;
}

.login-field span {
	font-size: 0.85rem;
	color: var(--p-text-muted-color);
}

.login-error {
	padding: 0.55rem 0.7rem;
	border-radius: 0.5rem;
	border: 1px solid color-mix(in oklab, var(--p-red-500) 40%, var(--p-content-border-color));
	background: color-mix(in oklab, var(--p-red-500) 10%, transparent);
	color: var(--p-red-500);
	font-size: 0.85rem;
}
</style>
