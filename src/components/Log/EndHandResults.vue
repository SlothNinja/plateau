<template>
  <GameEntry class='ma-1'>

    <template #toolbar-title>
      Hand: {{handNumber}} 
    </template>

    <Board :dTeam='dTeam' :numPlayers='numPlayers' :tricks='tricks' :path='path' />
    <Bids :bids='bids' :dTeam='dTeam' :order='order'/>
    <div>Declarers {{success}} bid.</div>

  </GameEntry>
</template>

<script setup>
import GameEntry from '@/snvue/components/Log/GameEntry'
import Bids from '@/components/Game/Bids'
import Board from '@/components/Game/Board'
import { gameKey } from '@/snvue/composables/keys.js'
import { computed, inject, unref } from 'vue'
import { bidValue } from '@/composables/bid.js'
import _get from 'lodash/get'
import _size from 'lodash/size'

const props = defineProps(['data'])
const game = inject(gameKey)

const bids = computed(() => _get(props, 'data.Results.Bids', []))
const dTeam = computed(() => _get(props, 'data.Results.DeclarersTeam', []))
const order = computed(() => _get(props, 'data.Results.SeatOrder', []))
const tricks = computed(() => _get(props, 'data.Results.Tricks', []))
const path = computed(() => _get(props, 'data.Results.Path', []))

const success = computed(() => {
  if (_get(props, 'data.Results.Success', '') == 'failure') {
    return 'failed'
  }
  return 'made'
})

const numPlayers = computed(() => _get(unref(game), 'Header.NumPlayers', 0))
const bid = computed(() => _get(props, 'data.Bid', {}))
const exchange = computed(() => _get(unref(bid), 'Exchange', ''))
const objective = computed(() => _get(unref(bid), 'Objective', ''))
const teams = computed(() => _get(unref(bid), '.Teams', ''))
const bValue = computed(() => bidValue(numPlayers, bid))

const handNumber = computed(() => _get(unref(props), 'data.Results.HandNumber', -1))
</script>
