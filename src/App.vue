<template>
  <router-view />
</template>

<script setup>
import { computed, provide, readonly, ref, unref, watch } from 'vue'
import { useFirebaseAuth } from 'vuefire'
import { signInWithCustomToken } from "firebase/auth";

/////////////////////////////////////////////////////
// get and provide current user
import  { useFetch } from '@/snvue/composables/fetch'
import { cuKey, credentialsKey } from '@/snvue/composables/keys'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

const cuURL = `${import.meta.env.VITE_PLATEAU_BACKEND}sn/user/current`
const { data, isFinished } = useFetch(cuURL).json()

const cu = computed(
  () => {
    if (unref(isFinished)) {
      let cuState = _get(unref(data), 'CU', null)
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
    if (unref(isFinished)) {
      return _get(unref(data), fsTokenKey, null)
    }
    return null
  }
)

const fsTokenKey = import.meta.env.VITE_FS_TOKEN_KEY
const hasFSToken = computed(() => !_isEmpty(fsTokenKey))
const credentials = ref({})
const auth = useFirebaseAuth()
watch(
  token,
  () => {
    signInWithCustomToken(auth, unref(token))
      .then((userCredential) => {
        credentials.value = userCredential
      })
      .catch((error) => {
        console.log(`errorCode: ${error.code} errorMessage: ${error.message}`)
      });
  }
)

provide( cuKey, readonly(cu) )
provide( credentialsKey, readonly(credentials) )
////////////////////////////////////////////////////
</script>
