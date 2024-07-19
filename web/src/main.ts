import { createApp } from 'vue'
import { createPinia } from 'pinia'

import router from './lib/router.ts'
import App from './App.vue'

import './assets/styles/tailwind.css'

const app = createApp(App)

app.use(router)
app.use(createPinia())

app.mount('#app')
