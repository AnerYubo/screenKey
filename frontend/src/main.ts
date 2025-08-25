import {createApp} from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'
import './style.css';
import naive from 'naive-ui'
import router from './router'
import { Buffer } from 'buffer'
window.Buffer = Buffer
const app = createApp(App)
const pinia = createPinia()
app.use(naive)
app.use(router)
app.use(pinia)
app.mount('#app')
document.addEventListener('contextmenu', e => e.preventDefault())