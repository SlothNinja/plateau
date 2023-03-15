import { ref } from 'vue'
import { get, includes } from 'lodash'

export function useGravatar(hash, size, t) {
  const gravSizes = { 'x-small': '24', 'small': '30', 'medium': '48', 'large': '54' }
  const gravTypes = [ 'personal', 'identicon', 'monsterid', 'retro', 'robohash' ]
  const sz = get(gravSizes, size, '64')

  if (t == 'personal') {
    return `https://www.gravatar.com/avatar/${hash}?s=${sz}`
  }

  if (!includes(gravTypes, t)) {
    t = 'monsterid'
  }
  return `https://www.gravatar.com/avatar/${hash}?s=${sz}&d=${t.value}&f=y`
}
