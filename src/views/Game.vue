<template>
  <v-container>
    <v-row dense>
      <v-col cols='6'>
        <v-row dense>
          <v-col cols='12'>
            <Info />
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <MessageBar />
          </v-col>
        </v-row>
        <v-row dense v-if='showTrick'>
          <v-col cols='12'>
            <Trick height='170' />
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <Hand class='pa-2' height='170' />
          </v-col>
        </v-row>
        <v-row dense v-if='showForm'>
          <v-col cols='12'>
            <BidForm />
          </v-col>
        </v-row>
      </v-col>
      <v-col cols='6'>
        <v-row dense>
          <v-col cols='12'>
            <Bids :bids='bids' :order='order' :dTeam='dTeam' />
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <v-radio-group inline v-model='results' >
              <v-radio v-for='hand in round' :label='label(hand)' :value='hand' :key='hand' />
            </v-radio-group>
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <Board :tricks='tricks' :declarersTeam='declarersTeam' />
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

// composables
import { cuKey, gameKey } from '@/composables/keys'
import { usePlayerByUser, useCP, useIsCP } from '@/composables/player'
import { useFetch } from '@/composables/fetch'

// vue
import { computed, inject, provide, ref, unref, watch} from 'vue'
// lodash
import _get from 'lodash/get'
import _size from 'lodash/size'
import _find from 'lodash/find'

const cu = inject(cuKey)
const game = inject(gameKey)

const player = computed(() => usePlayerByUser(game, cu))

const height = 170

const round = computed(() => _get(unref(game), 'Round', 1))

const results = ref(0)

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

const isCP = computed(() => (useIsCP(game, cu)))
const cp = computed(() => (useCP(game)))

const phase = computed(() => _get(unref(game), 'Phase', ''))
const showForm = computed(() => {
  return (unref(isCP) && !unref(cp).PerformedAction) && (unref(phase) == 'bid' || unref(phase) == 'increase objective')
})

const showTrick = computed(() => {
  return unref(phase) == 'card play'
})

</script>
