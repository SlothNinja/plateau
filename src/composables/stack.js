import { unref, computed } from 'vue'
import { usePlayerByPID } from '@/composables/player'
import _isEmpty from 'lodash/isEmpty'

function stacksFor(player) {
  const p = unref(player)
  if (_isEmpty(p)) {
    return []
  }
  return [ p.Stack0, p.Stack1, p.Stack2, p.Stack3, p.Stack4 ]
}

export function useStackByPID(game, pid) {
  const p = usePlayerByPID(game, pid)
  return stacksFor(p)
}
