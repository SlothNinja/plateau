<template>
  <v-container fluid>

    <v-row dense>
      <v-col cols='6'>
        <v-row dense>
          <v-col cols='12'>
            <Info :dTeam='dTeam'/>
          </v-col>
        </v-row>
        <v-row dense class='h-50'>
          <v-col cols='12'>
            <MessageBar/>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols='6'>
        <v-row dense class='h-100'>
          <v-col cols='12'>
            <v-card elevation='4' class='h-100'>
              <v-card-text>
                <Bids :bids='bids' :order='order' :dTeam='dTeam' />
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>

    <v-row dense>
      <v-col cols='6'>
        <v-row dense>
          <v-col cols='12'>
            <v-row dense v-if='showStack'>
              <v-col>
                <v-card elevation='4'>
                  <v-card-title :class='textcolor'>{{title}}</v-card-title>
                  <v-card-text class='h-100 w-100' >
                    <CardDisplay :height='height' :multi='1' :cards='oStacks'/>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
            <v-row dense v-if='showForm'>
              <v-col>
                <BidForm/>
              </v-col>
            </v-row>
          </v-col>
        </v-row>
        <v-row dense v-if='showTrick36'>
          <v-col>
            <Trick :height='height'/>
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <v-card>
              <Hand class='pa-2' :height='height' />
            </v-card>
          </v-col>
        </v-row>
      </v-col>
      <v-col cols='6'>
        <v-row dense v-if='showTrick2'>
          <v-col>
            <Trick :height='height'/>
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <v-card>
            <Board :tricks='tricks' :dTeam='dTeam' :numPlayers='numPlayers' :path='path' />
            </v-card>
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <v-card class='h-100'>
              <v-card-text>
                <v-radio-group hide-details inline v-model='results' >
                  <v-radio v-for='hand in round' :label='label(hand)' :value='hand' :key='hand' />
                </v-radio-group>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
// components
import Bids from '@/components/Game/Bids'
import BidForm from '@/components/Game/BidForm'
import Trick from '@/components/Game/Trick'
import Hand from '@/components/Game/Hand'
import Info from '@/components/Game/Info'
import Board from '@/components/Game/Board'
import MessageBar from '@/components/Game/MessageBar'
import CardDisplay from '@/components/Game/CardDisplay'

// composables
import { cuKey, gameKey } from '@/snvue/composables/keys'
import { usePlayerByUser, usePlayerByPID, useNameFor, useCP, useCPID, useIsCP } from '@/composables/player'
import { useFetch } from '@/snvue/composables/fetch'
import { useStackByPID } from '@/composables/stack'
import { useColorFor } from '@/composables/color'

// vue
import { computed, inject, provide, ref, unref, watch} from 'vue'
// lodash
import _get from 'lodash/get'
import _size from 'lodash/size'
import _find from 'lodash/find'

const cu = inject(cuKey)
const game = inject(gameKey)

const header = computed(() => _get(unref(game), 'Header', {}))
const player = computed(() => usePlayerByUser(game, cu))
const pid = computed(() => _get(unref(player), 'ID', 1))
const opid = computed(
  () => {
    switch (unref(pid)) {
      case 1:
        return 2
      case 2:
        return 1
      default:
        return 0
    }
  }
)

const height = 170

const round = computed(() => _get(unref(header), 'Round', 1))

const results = ref(0)

watch(round, () => (results.value = unref(round)))

const oStacks = computed(() => useStackByPID(game, opid))
const title = computed(() => unref(useNameFor(header, opid)))

const index = computed(
  () => {
    if (unref(results) == 0) {
      results.value = unref(round)
    }
    return unref(results) - 1
  }
)

const tricks = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(game), 'State.Tricks', [])
    }
    return _get(unref(game), `State.LastResults[${unref(index)}].Tricks`, [])
  }
)

const path = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return []
    }
    return _get(unref(game), `State.LastResults[${unref(index)}].Path`, [])
  }
)

const dTeam = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(game), 'State.DeclarersTeam', [])
    }
    return _get(unref(game), `State.LastResults[${unref(index)}].DeclarersTeam`, [])
  }
)

const numPlayers = computed(() => _get(unref(header), 'NumPlayers', 0))

const showStack = computed(() => unref(numPlayers) == 2)

const bids = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(game), 'State.Bids', [])
    }
    return _get(unref(game), `State.LastResults[${unref(index)}].Bids`, [])
  }
)

const order = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(header), 'OrderIDS', [])
    }
    return _get(unref(game), `State.LastResults[${unref(index)}].SeatOrder`, [])
  }
)

const textcolor = computed(() => `text-${useColorFor(dTeam, opid)}`)
          
function label(hand) {
  if (hand == unref(round)) {
    return 'Current'
  }
  return `Hand ${unref(hand)}`
}

const phase = computed(() => _get(unref(header), 'Phase', ''))
const showForm = computed(() => {
  return (unref(phase) == 'bid' || unref(phase) == 'increase objective')
})

const showTrick2 = computed(() => {
  return unref(numPlayers) == 2 && unref(phase) == 'card play'
})

const showTrick36 = computed(() => {
  return unref(numPlayers) != 2 && unref(phase) == 'card play'
})

</script>
