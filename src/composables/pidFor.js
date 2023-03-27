import _get from 'lodash/get'
import _findIndex from 'lodash/findIndex'

export function usePIDFor(game, user) {
  const uid = _get(user, 'value.id', -1)
  const uids = _get(game, 'value.header.userIds', []) 
  const index = _findIndex(uids, id => id == uid)
  return index + 1
}
