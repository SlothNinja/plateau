import { useUser } from '@/composables/user.js'
import { map } from 'lodash'

export function useUsers(header) {
  return map(header.userIds, function (id, i) {
    return useUser(header, i)
  })
}
