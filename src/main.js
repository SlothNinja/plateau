// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'
import { VueFire, VueFireAuth } from 'vuefire'
import { firebaseApp } from '@/composables/firebase'

// Plugins
import { registerPlugins } from '@/plugins'

const app = createApp(App)

// creat app
app
  .use(VueFire, {
    // imported above but could also just be created here
    firebaseApp,
    modules: [
      // we will see other modules later on
      VueFireAuth(),
    ],
  })

registerPlugins(app)

app.mount('#app')
