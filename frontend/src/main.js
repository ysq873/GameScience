import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createStore } from 'vuex'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import routes from './router'
import storeConfig from './store'
import { getSession } from './api/auth'

const app = createApp(App)

const router = createRouter({
  history: createWebHistory(),
  routes
})

const store = createStore(storeConfig)

app.use(router)
app.use(store)
app.use(ElementPlus)

Promise.resolve()
  .then(() => getSession())
  .then(me => {
    const id = me?.data?.identity?.id
    if (id) {
      const traits = me.data.identity.traits || {}
      store.commit('SET_USER', {
        id,
        email: traits.email || '',
        name: traits.name || '',
        session: me.data
      })
    }
  })
  .finally(() => {
    app.mount('#app')
  })