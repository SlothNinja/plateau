<template>
  <v-container fluid>
    <v-row dense>
      <v-col cols='6'>
        <v-row dense class='h-50'>
          <v-col cols='12'>
            <Info />
          </v-col>
        </v-row>
        <v-row dense class='h-50'>
          <v-col cols='12'>
            <MessageBar />
          </v-col>
        </v-row>
      </v-col>
      <v-col cols='6'>
        <v-row dense class='h-100'>
          <v-col cols='12'>
            <Bids :bids='bids' :order='order' :dTeam='dTeam' />
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
                <StacksDisplay :height='height' :title='title' :stacks='oStacks' />
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
            <Board :tricks='tricks' :declarersTeam='declarersTeam' :numPlayers='numPlayers' />
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
import Bids from '@/components/Game/Bids.vue'
import BidForm from '@/components/Game/BidForm.vue'
import Trick from '@/components/Game/Trick.vue'
import Hand from '@/components/Game/Hand.vue'
import Info from '@/components/Game/Info.vue'
import Board from '@/components/Game/Board.vue'
import Table from '@/components/Game/Table.vue'
import MessageBar from '@/components/Game/MessageBar.vue'
import StacksDisplay from '@/components/Game/StacksDisplay.vue'

// composables
import { cuKey, gameKey } from '@/composables/keys'
import { usePlayerByUser, usePlayerByPID, useNameFor, useCP, useCPID, useIsCP } from '@/composables/player'
import { useFetch } from '@/composables/fetch'
import { useStackByPID } from '@/composables/stack'

// vue
import { computed, inject, provide, ref, unref, watch} from 'vue'
// lodash
import _get from 'lodash/get'
import _size from 'lodash/size'
import _find from 'lodash/find'

const cu = inject(cuKey)
const game = inject(gameKey)

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

const round = computed(() => _get(unref(game), 'Round', 1))

const results = ref(0)

watch(round, () => (results.value = unref(round)))

const oStacks = computed(() => useStackByPID(game, opid))
const title = computed(() => useNameFor(game, opid))

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
      return _get(unref(game), 'Tricks', [])
    }
    return _get(unref(game), `LastResults[${unref(index)}].Tricks`, [])
  }
)

const declarersTeam = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(game), 'DeclarersTeam', [])
    }
    return _get(unref(game), `LastResults[${unref(index)}].DeclarersTeam`, [])
  }
)

const numPlayers = computed(() => _get(unref(game), 'NumPlayers', 0))

const showStack = computed(() => unref(numPlayers) == 2)

const bids = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(game), 'Bids', [])
    }
    return _get(unref(game), `LastResults[${unref(index)}].Bids`, [])
  }
)

const order = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(game), 'OrderIDS', [])
    }
    return _get(unref(game), `LastResults[${unref(index)}].SeatOrder`, [])
  }
)

const dTeam = computed(
  () => {
    if ((unref(results) == unref(round))) {
      return _get(unref(game), 'DeclarersTeam', [])
    }
    return _get(unref(game), `LastResults[${unref(index)}].DeclarersTeam`, [])
  }
)

function label(hand) {
  if (hand == unref(round)) {
    return 'Current'
  }
  return `Hand ${unref(hand)}`
}

const phase = computed(() => _get(unref(game), 'Phase', ''))
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
