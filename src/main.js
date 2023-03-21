/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Components
import App from './App.vue'

// Composables
import { createApp, computed, readonly } from 'vue'

// Plugins
import { registerPlugins } from '@/plugins'

// creat app
const app = createApp(App)

/////////////////////////////////////////////////////
// get and provide current user
import  { useFetch } from '@/composables/fetch.js'
import { get } from 'lodash'
import { cuKey } from '@/composables/keys.js'

const { data, error } = useFetch('/sn/cu')
const cu = computed( () => {
 return get(data, 'value.cu', {})
})

app.provide( cuKey, readonly(cu) )
////////////////////////////////////////////////////

registerPlugins(app)

app.mount('#app')
