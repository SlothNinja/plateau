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
                  :user='useUserByIndex(game.header, p.id-1)'
                  :size='24'
                  :color='colorFor(p.id)'
                  />
            </td>
            <td class='text-center'>{{bidLabel(p.id)}}</td>
            <td class='text-center'>{{p.score}}</td>
          </tr>
        </tbody>
      </v-table>
    </v-card-text>
  </v-card>
</template>

<script setup>
// import components
import UserButton from '@/components/UserButton.vue'

// import composables
import { bidValue } from '@/composables/bid.js'
import { useUserByIndex } from '@/composables/user.js'
import { useIsCP, usePlayerByPID } from '@/composables/player.js'

// import lodash
import _find from 'lodash/find'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import _includes from 'lodash/includes'

// import vue
import { computed, inject } from 'vue'

const players = computed(() => _get(game, 'value.state.players', []))

function bidLabel(pid) {
  const p = usePlayerByPID(game, pid)
  if (p.passed) {
    return 'passed'
  }
  const bid = _find(_get(game, 'value.state.bids', []), [ 'pid', pid ])
  if (_isEmpty(bid)) {
    return 'no bid'
  }
  const exchange = _get(bid, 'exchange', '')
  const objective = _get(bid, 'objective', '')
  const teams = _get(bid, 'teams', '')
  const bValue = bidValue(game, bid)
  return `${exchange} ${objective} ${teams} (${bValue})`
}

// inject game and current user
import { cuKey, gameKey } from '@/composables/keys.js'
const cu = inject(cuKey)
const { game, updateGame } = inject(gameKey)

function cpClass(player) {
  const pid = player.id

  if (_includes(_get(game, 'value.header.cpids', []), pid)) {
    if (useIsCP(game, cu)) {
      return 'font-weight-black text-red-darken-4'
    }
    return 'font-weight-black'
  }
  return ''
}

function colorFor(pid) {
  if (_includes(_get(game, 'value.state.declarersTeam', []), pid)) {
    return 'rgb(150 0 0)'
  }
  return 'rgb(0 0 150)'
}
</script>
