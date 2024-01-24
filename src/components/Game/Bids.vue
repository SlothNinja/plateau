<template>
  <v-card elevation='4' class='h-100'>
    <v-card-text>
      <v-table density='compact'>
        <thead>
          <tr>
            <th class='text-center'>Name</th>
            <th class='text-center'>Passed</th>
            <th class='text-center'>Last Bid</th>
            <th class='text-center'>Score</th>
          </tr>
        </thead>
        <tbody>
          <tr :class='cpClass(p)' v-for='(p, index) in players' :key='index'>
            <td class='text-center'>
              <UserButton
                  :user='useUserByPID(header, p.ID)'
                  :size='24'
                  :color='useColorFor(dTeam, p.ID)'
                  >
                  {{nameFor(p)}}
              </UserButton>
            </td>
            <td class='text-center'>{{p.Passed}}</td>
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
import { bidValue } from '@/composables/bid'
import { useUserByPID } from '@/composables/user'
import { useIsCP, useCPID, useNameFor, usePlayerByPID } from '@/composables/player'
import { useColorFor } from '@/composables/color'

// import lodash
import _findLast from 'lodash/findLast'
import _first from 'lodash/first'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import _includes from 'lodash/includes'
import _map from 'lodash/map'

// import vue
import { computed, inject, unref } from 'vue'

const props = defineProps(['bids', 'order', 'dTeam'])

const players = computed(() => _map(_get(props, 'order', []), (pid) => usePlayerByPID(game, pid)))

function bidLabel(pid) {
  const bid = _findLast(_get(props, 'bids', []), [ 'PID', pid ])
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

const header = computed(() => _get(unref(game), 'Header', {}))

function declarer(player) {
  return _first(_get(unref(game), 'State.DeclarersTeam', [])) == unref(player).ID
}

const cpid = computed(() => useCPID(header))

function nameFor(p) {
  let name = useNameFor(header, p.ID)
  if (declarer(p)) {
    return `[${name}]`
  }
  return name
}

function cpClass(player) {
  const pid = player.ID
  const color = useColorFor(props.dTeam, pid)
 
  if (_includes(_get(unref(header), 'CPIDS', []), pid)) {
    return `font-weight-black text-${color}`
  }
  return `text-${color}`
}
</script>
