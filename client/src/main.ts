import { createApp } from 'vue'
import {createPinia} from 'pinia'
import App from './App.vue'
import router from './router/router.ts'
import './index.css'
import Toast from 'vue-toastification'
import 'vue-toastification/dist/index.css'
import 'leaflet/dist/leaflet.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(Toast)

app.mount('#app')