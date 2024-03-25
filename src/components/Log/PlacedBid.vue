<template>
  <PlayerEntry class='ma-1' :pid='pid'>
    <template #toolbar-title>
      Hand: {{handNumber}}
    </template>
    <div v-if='teams'>
      Placed a bid of "{{exchange}} {{objective}} {{teams}}" for a value of {{bValue}}.
    </div>
    <div v-else>
      Placed a bid of "{{exchange}} {{objective}}" for a value of {{bValue}}.
    </div>
  </PlayerEntry>
</template>

<script setup>
import PlayerEntry from '@/snvue/components/Log/PlayerEntry'
import { gameKey } from '@/snvue/composables/keys.js'
import { computed, inject, unref } from 'vue'
import _get from 'lodash/get'
import { bidValue } from '@/composables/bid.js'
const props = defineProps(['data'])

const game = inject(gameKey)
const bid = computed(() => _get(props, 'data.Bid', {}))
const exchange = computed(() => _get(unref(bid), 'Exchange', ''))
const objective = computed(() => _get(unref(bid), 'Objective', ''))
const teams = computed(() => _get(unref(bid), '.Teams', ''))
const numPlayers = computed(() => _get(unref(game), 'Header.NumPlayers', 0))
const bValue = computed(() => bidValue(numPlayers, bid))

const pid = computed(() => _get(unref(props), 'data.PID', ''))
const handNumber = computed(() => _get(unref(props), 'data.HandNumber', -1))
</script>
