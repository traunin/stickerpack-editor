import { createApp } from 'vue'
import router from "@/routes.ts"

import '@/assets/styles/reset.css'
import '@/assets/styles/style.css'

import App from '@/App.vue'

const app = createApp(App)

app.use(router)

app.mount('#app')
