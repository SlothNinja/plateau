<template>
  <v-sheet elevation='4' rounded width='90%' class='pa-4'>
    <v-container>
    <v-row dense>
      <v-col cols='12' class='text-center'>
        Set Your Bid Characteristics
      </v-col>
    </v-row>
    <v-row dense>
      <v-col cols='9'>
      </v-col>
      <v-col cols='3' class='text-center'>
        Point Value
      </v-col>
    </v-row>
    <v-row dense v-if='showExchange'>
      <v-col cols='4'>
        Card Exchange
      </v-col>
      <v-col cols='5'>
        <v-radio-group v-model='bid.exchange'>
          <v-radio :label="'Exchange (' + exValue('exchange') + ')'" value='exchange'/>
          <v-radio :label="'No Exchange (' + exValue('no-exchange') + ')'" value='no-exchange'/>
        </v-radio-group>
      </v-col>
      <v-col cols='3' class='d-flex justify-center align-center'>
        <div>{{exchangeValue}}</div>
      </v-col>
    </v-row>
    <v-row dense>
      <v-col cols='4'>
        Objective
      </v-col>
      <v-col cols='5'>
        <v-radio-group v-model='bid.objective'>
          <v-radio :label="'Bridge (' + obValue('bridge') + ')'" value='bridge'/>
          <v-radio :label="'Y (' + obValue('y') + ')'" value='y'/>
          <v-radio :label="'Fork (' + obValue('fork') + ')'" value='fork'/>
          <v-radio :label="'5 sides (' + obValue('5-sides') + ')'" value='5-sides'/>
          <v-radio :label="'6 sides (' + obValue('6-sides') + ')'" value='6-sides'/>
        </v-radio-group>
      </v-col>
      <v-col cols='3' class='d-flex justify-center align-center'>
        <div>{{objectiveValue}}</div>
      </v-col>
    </v-row>
    <v-row dense v-if='showTeams'>
      <v-col cols='4'>
        Teams
      </v-col>
      <v-col cols='5'>
        <v-radio-group v-model='bid.teams'>
          <v-radio v-if='showTrio' :label="'Trio (' + tValue('trio') + ')'" value='trio'/>
          <v-radio :label="'Duo (' + tValue('duo') + ')'" value='duo'/>
          <v-radio :label="'Solo (' + tValue('solo') + ')'" value='solo'/>
        </v-radio-group>
      </v-col>
      <v-col cols='3' class='d-flex justify-center align-center'>
        <div>{{teamsValue}}</div>
      </v-col>
    </v-row>
    <v-row dense>
      <v-col cols='4'>
        Total Bid Value:
      </v-col>
      <v-col cols='5'>
      </v-col>
      <v-col cols='3' class='d-flex justify-center align-center'>
        <div>{{bidValue}}</div>
      </v-col>
    </v-row>
    <v-row v-if='!cp.performedAction' class='d-flex justify-center align-center'>
      <v-btn color='green' size='small' @click='submit'>Submit</v-btn>
    </v-row>
  </v-container>
  </v-sheet>
</template>

<script setup>
import { computed, inject, ref, watch, onMounted } from 'vue'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

import { usePut } from '@/composables/put.js'

// inject game and current user
import { cuKey, gameKey, snackKey } from '@/composables/keys.js'
const cu = inject(cuKey)
const { game, updateGame } = inject(gameKey)

import { usePlayerFor } from '@/composables/playerFor.js'
const cp = computed(() => usePlayerFor(game, cu))
const cpid = computed(() => _get(cp, 'value.id', -1))

const bvalue = ref({})

const bid = computed({
  get() {
    if (_isEmpty(bvalue.value)) {
      bvalue.value = minBid.value
    }
    return bvalue.value
  },
  set(value) {
    bvalue.value = value
  }
})

const minBid = computed(() => {
  switch (numPlayers.value) {
    case 2:
      return { exchange: 'exchange', objective: 'y', pid: cpid.value }
    case 3:
      return { exchange: 'exchange', objective: 'bridge', pid: cpid.value }
    case 4:
      return { exchange: 'exchange', objective: 'y', teams: 'duo', pid: cpid.value }
    case 5:
      return { exchange: 'exchange', objective: 'bridge', teams: 'duo', pid: cpid.value }
    case 6:
      return { objective: 'y', teams: 'trio', pid: cpid.value }
    default:
      return { exchange: '', objective: '', teams: '', pid: -1 }
  }
})

const exchangeValue = computed(() => exValue(bid.value.exchange))

function exValue(bid) {
  return (bid == 'exchange') ? 1 : 2
}

const objectiveValue = computed(() => obValue(bid.value.objective))

function  obValue(bid) {
  switch (bid) {
    case 'bridge':
      return 0
    case 'y':
      return 2
    case 'fork':
      return 4
    case '5-sides':
      return 6
    case '6-sides':
      return 8
    default:
      return 0
  }
}

const numPlayers = computed(() => _get(game, 'value.header.numPlayers', 0))
const showExchange = computed(() => (numPlayers.value < 6))
const showTeams = computed(() => (numPlayers.value >= 4))
const showTrio = computed(() => (numPlayers.value == 6))

const teamsValue = computed(() => tValue(bid.value.teams))

// watch(numPlayers, () => (bid.value = minBid.value))

function tValue(bid) {
  switch (numPlayers.value) {
    case 4:
      return team45(bid)
    case 5:
      return team45(bid)
    case 6:
      return team6(bid)
    default:
      return 0
  }
}

function team45(bid) {
  return (bid == 'solo') ? 5 : 0
}

function team6(bid) {
  switch (bid) {
    case 'duo':
      return 5
    case 'solo':
      return 10
    default:
      return 0
  }
}

const bidValue = computed(() => (exchangeValue.value + objectiveValue.value + teamsValue.value))

//////////////////////////////////////
// Snackbar
const snackbar = inject(snackKey)

function submit () {
  const { response, error } = usePut(`/sn/game/bid/${game.value.id}`, bid)

  watch( response, () => {
    const g = _get(response, 'value.game', {})
    if (!_isEmpty(g)) {
      updateGame(g)
    }
    const msg = _get(response, 'value.message', '')
    if (!_isEmpty(msg)) {
      snackbar.value.message = msg
      snackbar.value.open = true
    }
  })
}

</script>
