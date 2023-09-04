// Components
import App from './App.vue'

// Composables
import { createApp, computed, readonly, ref, unref, watch } from 'vue'
import { VueFire, VueFireAuth, useFirebaseAuth } from 'vuefire'
import { firebaseApp } from '@/composables/firebase'
import { signInWithCustomToken } from "firebase/auth";

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
import  { useFetch } from '@/composables/fetch'
import _get from 'lodash/get'
import { cuKey, credentialsKey } from '@/composables/keys'
import _isEmpty from 'lodash/isEmpty'

const fetchURL = `${import.meta.env.VITE_PLATEAU_BACKEND}sn/user/current`
const { state, isLoading, isReady } = useFetch(fetchURL)
const cu = computed(
  () => {
    if (unref(isReady)) {
      let cuState = _get(unref(state), 'CU', null)
      if (unref(hasFSToken)) {
        return _isEmpty(unref(credentials)) ? null : cuState
      }
      return cuState
    }
    return null
  }
)

const token = computed(
  () => {
    if (unref(isReady)) {
      return _get(unref(state), fsTokenKey, null)
    }
    return null
  }
)

const fsTokenKey = import.meta.env.VITE_FS_TOKEN_KEY
const hasFSToken = computed(() => !_isEmpty(fsTokenKey))
const credentials = ref({})
watch(
  token,
  () => {
    signInWithCustomToken(useFirebaseAuth(), unref(token))
      .then((userCredential) => {
        credentials.value = userCredential
      })
      .catch((error) => {
        console.log(`errorCode: ${error.code} errorMessage: ${error.message}`)
      });
  }
)

app.provide( cuKey, readonly(cu) )
app.provide( credentialsKey, readonly(credentials) )
////////////////////////////////////////////////////

registerPlugins(app)

app.mount('#app')
