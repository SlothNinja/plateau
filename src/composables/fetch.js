import { ref, unref } from 'vue'

export function useFetch(url) {
  const data = ref(null)
  const error = ref(null)

  if (process.env.NODE_ENV == 'development') {
    url = 'https://plateau.fake-slothninja.com:8091' + url
    console.log('fetching: ' + unref(url))
  }

  fetch(unref(url), { credentials: 'include' } )
    .then((res) => res.json())
    .then((json) => (data.value = json))
    .catch((err) => (error.value = err))

  return { data, error }
}
