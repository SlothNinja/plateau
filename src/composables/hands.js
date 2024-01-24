import { unref } from 'vue'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

function handsPerPlayer(game) {
  const opt = _get(unref(game), 'Header.OptString', {})
  if (_isEmpty(opt)) {
    return 0
  }
  return _get(JSON.parse(opt), 'HandsPerPlayer', 0)
}

export function useHands(game) {
  return handsPerPlayer(game) * _get(unref(game), 'Header.NumPlayers', 0)
}
