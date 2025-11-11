import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createStore } from 'vuex'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import routes from './router'
import storeConfig from './store'

const app = createApp(App)

const router = createRouter({
  history: createWebHistory(),
  routes
})

const store = createStore(storeConfig)

app.use(router)
app.use(store)
app.use(ElementPlus)

app.mount('#app')