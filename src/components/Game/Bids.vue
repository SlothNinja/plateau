<template>
  <v-card elevation='4'>
    <v-card-text>
      <v-table density='compact'>
        <thead>
          <tr>
            <th class='text-center'>Dealer</th>
            <th class='text-center'>Name</th>
            <th class='text-center'>Last Bid</th>
            <th class='text-center'>Score</th>
          </tr>
        </thead>
        <tbody>
          <tr :class='cpClass(p)' v-for='(p, index) in players' :key='index'>
            <td class='text-center'>{{(index == 0) ? 'X' : ''}}</td>
            <td class='text-center'>
              <UserButton
                  :user='useUserByIndex(game, p.ID-1)'
                  :size='24'
                  :color='useColorFor(dTeam, p.ID)'
                  />
            </td>
            <td class='text-center'>{{bidLabel(p.ID)}}</td>
            <td class='text-center'>{{p.Score}}</td>
          </tr>
        </tbody>
      </v-table>
    </v-card-text>
  </v-card>
</template>

<script setup>
// import components
import UserButton from '@/components/Common/UserButton.vue'

// import composables
import { bidValue } from '@/composables/bid.js'
import { useUserByIndex } from '@/composables/user.js'
import { useIsCP, usePlayerByPID } from '@/composables/player.js'
import { useColorFor } from '@/composables/color.js'

// import lodash
import _find from 'lodash/find'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import _includes from 'lodash/includes'
import _map from 'lodash/map'

// import vue
import { computed, inject, unref } from 'vue'

const props = defineProps(['bids', 'order', 'dTeam'])

const players = computed(() => _map(_get(props, 'order', []), (pid) => usePlayerByPID(game, pid)))

function bidLabel(pid) {
  const p = usePlayerByPID(game, pid)
  if (p.Passed) {
    return 'passed'
  }
  const bid = _find(_get(props, 'bids', []), [ 'PID', pid ])
  if (_isEmpty(bid)) {
    return 'no bid'
  }
  const exchange = _get(bid, 'Exchange', '')
  const objective = _get(bid, 'Objective', '')
  const teams = _get(bid, 'Teams', '')
  const bValue = bidValue(game, bid)
  return `${exchange} ${objective} ${teams} (${bValue})`
}

// inject game and current user
import { cuKey, gameKey } from '@/composables/keys.js'
const cu = inject(cuKey)
const game = inject(gameKey)

function cpClass(player) {
  const pid = player.id

  if (_includes(_get(unref(game), 'CPIDS', []), pid)) {
    if (useIsCP(game, cu)) {
      return 'font-weight-black text-red-darken-4'
    }
    return 'font-weight-black'
  }
  return ''
}
</script>
