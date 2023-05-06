// Components
import App from './App.vue'

// Composables
import { createApp, computed, readonly, unref } from 'vue'
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

/////////////////////////////////////////////////////
// get and provide current user
import  { useFetch } from '@/composables/fetch.js'
import { get } from 'lodash'
import { cuKey } from '@/composables/keys.js'

const { data, error } = useFetch('/sn/cu')
const cu = computed( () => {
  return get(unref(data), 'CU', {})
})

app.provide( cuKey, readonly(cu) )
////////////////////////////////////////////////////

registerPlugins(app)

app.mount('#app')
