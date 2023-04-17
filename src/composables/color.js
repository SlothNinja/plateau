import { unref } from 'vue'
import _get from 'lodash/get'
import _includes from 'lodash/includes'

export function useColorFor(game, pid) {
  if (_includes(_get(unref(game), 'value.state.declarersTeam', []), unref(pid))) {
    return 'rgb(150 0 0)'
  }
  return 'rgb(0 0 150)'
}
