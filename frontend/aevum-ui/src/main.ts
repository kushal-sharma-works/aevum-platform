import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import { installPrimeVue } from './plugins/primevue'
import router from './router'
import 'primeicons/primeicons.css'
import './styles/tailwind.css'
import './styles/main.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
installPrimeVue(app)

window.addEventListener('aevum:auth-required', async () => {
	if (router.currentRoute.value.path !== '/login') {
		await router.replace('/login')
	}
})

app.mount('#app')
