import { ref, unref } from 'vue'
import { useAsyncState } from '@vueuse/core'

export function useFetch(url) {
  if (process.env.NODE_ENV == 'development') {
    console.log('fetching: ' + unref(url))
  }

  return useAsyncState(fetch(unref(url), { credentials: 'include' } ).then((res) => res.json()), {})
}

function put(url, data) {
  if (process.env.NODE_ENV == 'development') {
    console.log('putting: ' + unref(url))
  }

  const opt = {
    method: 'PUT',
    credentials: 'include',
    headers: {
      'Content-type': 'application/json'
    },
    body: JSON.stringify(unref(data)),
  } 

  return fetch(unref(url), opt).then((res) => res.json())
}

export function usePut(url, data) {
  const response = ref(null)
  const error = ref(null)

  if (process.env.NODE_ENV == 'development') {
    console.log('putting: ' + unref(url))
  }

  return useAsyncState(put(url, data, {}))
}
