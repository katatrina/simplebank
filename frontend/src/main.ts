import './assets/main.css' // global css

import { createApp } from 'vue' // imported function APIs
import App from './App.vue' // import root component
import router from './router' // import router
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'

const app = createApp(App) // create application instance
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      cssLayer: {
        name: 'primevue',
        order: 'tailwind-base, primevue, tailwind-utilities',
      },
    },
  },
})

app.use(router) // register router

app.mount('#app') // call mount() to render UI, it return the root component "App"
