import { useUser } from '@/composables/user.js'
import { unref } from 'vue'
import _map from 'lodash/map'
import _get from 'lodash/get'

export function useUsers(header) {
  const h = unref(header)
  return _map(_get(h, 'userIds', []), (id, i) => useUser(h, i))
}
