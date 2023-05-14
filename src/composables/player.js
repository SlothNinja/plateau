import { unref } from 'vue'
import _find from 'lodash/find'
import _get from 'lodash/get'
import _findIndex from 'lodash/findIndex'
import { useUserByIndex } from '@/composables/user'

export function usePlayerByUser(game, user) {
  const pid = usePIDForUser(game, user)
  return usePlayerByPID(game, pid)
}

export function usePlayerByPID(game, pid) {
  const players = _get(unref(game), 'Players', [])
  return _find(players, [ 'ID', unref(pid) ])
}

export function usePIDForUser(game, user) {
  const g = unref(game)
  const u = unref(user)
  const uid = _get(u, 'ID', -1)
  const uids = _get(g, 'UserIDS', []) 
  const index = _findIndex(uids, id => id == uid)
  return index + 1
}

export function useCP(game) {
  const cpid = useCPID(game)
  const players = _get(unref(game), 'Players', [])
  return _find(players, [ 'ID', cpid ])
}

export function useCPID(game) {
  const g = unref(game)
  return _get(g, 'CPIDS[0]', -1)
}

export function useIsCP(game, cu) {
  return usePIDForUser(game, cu) == useCPID(game)
}

export function useNameFor(game, pid) {
  const user = useUserByIndex(unref(game), unref(pid) - 1)
  return _get(user, 'Name', '')
}
