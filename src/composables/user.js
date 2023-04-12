import { unref } from 'vue'
import _get from 'lodash/get'
import _map from 'lodash/map'

export function useUserByIndex(header, index) {
  const h = unref(header)
  const i = unref(index)
  return {
    id: _get(h, `userIds[${i}]`, 0),
    name: _get(h, `userNames[${i}]`, ''),
    emailHash: _get(h, `userEmailHashes[${i}]`, ''),
    gravType: _get(h, `userGravTypes[${i}]`, ''),
  }
}

export function useUsers(header) {
  const h = unref(header)
  return _map(_get(h, 'userIds', []), (id, i) => useUserByIndex(h, i))
}

export function useCreator(header) {
  const h = unref(header)
  return {
    id: h.creatorId,
    name: h.creatorName,
    emailHash: h.creatorEmailHash,
    gravType: h.creatorGravType
  }
}
