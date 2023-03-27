import { usePIDFor } from '@/composables/pidFor.js'
import { unref } from 'vue'
import _find from 'lodash/find'
import _get from 'lodash/get'

export function usePlayerFor(game, user) {
  const pid = usePIDFor(game, user)
  const players = _get(unref(game), 'state.players', [])
  return _find(players, [ 'id', pid ])
}
