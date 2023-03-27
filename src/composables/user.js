import { unref } from 'vue'
import _get from 'lodash/get'

export function useUser(header, index) {
  const h = unref(header)
  const i = unref(index)
  return {
    id: _get(h, `userIds[${i}]`, 0),
    name: _get(h, `userNames[${i}]`, ''),
    emailHash: _get(h, `userEmailHashes[${i}]`, ''),
    gravType: _get(h, `userGravTypes[${i}]`, ''),
  }
}
