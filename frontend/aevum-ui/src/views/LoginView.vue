<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
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
	<div class="row items-center justify-center" style="min-height: 100vh">
		<div class="col-11 col-sm-8 col-md-5 col-lg-4">
			<q-card flat bordered>
				<q-card-section>
					<div class="text-h6 text-weight-medium">Aevum Login</div>
					<div class="text-caption text-grey-7">Sign in to continue</div>
				</q-card-section>
				<q-separator />
				<q-card-section>
					<q-form class="q-gutter-sm" @submit.prevent="submit">
						<q-input v-model="username" label="Username" outlined />
						<q-input v-model="password" label="Password" type="password" outlined />
						<q-btn type="submit" color="primary" label="Login" class="full-width" />
					</q-form>
				</q-card-section>
			</q-card>
		</div>
	</div>
</template>
