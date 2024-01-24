import { unref } from 'vue'
import _get from 'lodash/get'
import _includes from 'lodash/includes'
import _isEmpty from 'lodash/isEmpty'

export function useColorFor(dTeam, pid) {
  if (_isEmpty(unref(dTeam))) {
    return 'black'
  }

  if (_includes(unref(dTeam), unref(pid))) {
    return 'red-accent-4'
  }
  return 'indigo-darken-4'
}
