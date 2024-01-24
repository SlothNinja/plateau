import { ref, unref } from 'vue'
import _find from 'lodash/find'
import _get from 'lodash/get'
import _findIndex from 'lodash/findIndex'
import { useUserByPID } from '@/composables/user'

export function usePlayerByUser(game, user) {
  const header = ref(_get(unref(game), 'Header', {}))
  const pid = usePIDForUser(header, user)
  return usePlayerByPID(game, pid)
}

export function usePlayerByPID(game, pid) {
  const players = _get(unref(game), 'Players', [])
  return _find(players, [ 'ID', unref(pid) ])
}

export function usePIDForUser(header, user) {
  const uid = _get(unref(user), 'ID', -1)
  const uids = _get(unref(header), 'UserIDS', []) 
  const index = _findIndex(uids, id => id == uid)
  return index + 1
}

export function useCP(game) {
  const header = _get(unref(game), 'Header', {})
  const cpid = useCPID(header)
  const players = _get(unref(game), 'Players', [])
  return _find(players, [ 'ID', cpid ])
}

export function useCPID(header) {
  return _get(unref(header), 'CPIDS[0]', -1)
}

export function useIsCP(header, cu) {
  return usePIDForUser(header, cu) == useCPID(header)
}

export function useNameFor(header, pid) {
  const user = useUserByPID(header, pid)
  return _get(user, 'Name', '')
}
