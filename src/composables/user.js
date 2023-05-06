import { unref } from 'vue'
import _get from 'lodash/get'
import _map from 'lodash/map'

export function useUserByIndex(header, index) {
  const h = unref(header)
  const i = unref(index)
  return {
    ID: _get(h, `UserIDS[${i}]`, 0),
    Name: _get(h, `UserNames[${i}]`, ''),
    EmailHash: _get(h, `UserEmailHashes[${i}]`, ''),
    GravType: _get(h, `UserGravTypes[${i}]`, ''),
  }
}

export function useUsers(header) {
  const h = unref(header)
  return _map(_get(h, 'UserIDS', []), (id, i) => useUserByIndex(h, i))
}

export function useCreator(header) {
  const h = unref(header)
  return {
    ID: h.CreatorId,
    Name: h.CreatorName,
    EmailHash: h.CreatorEmailHash,
    GravType: h.CreatorGravType
  }
}
