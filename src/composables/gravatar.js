import { ref } from 'vue'
import { get, includes } from 'lodash'

export function useGravatar(hash, size, t) {
  const gravSizes = { 'x-small': '24', 'small': '30', 'medium': '48', 'large': '54' }
  const gravTypes = [ 'personal', 'identicon', 'monsterid', 'retro', 'robohash' ]
  const path = ref('')
  const sz = get(gravSizes, size, '64')

  if (!includes(gravTypes, t)) {
    t = 'monsterid'
  }
  path.value = `https://www.gravatar.com/avatar/${hash}?s=${sz}&d=${t}&f=y`
  return path
}
