<template>
  <v-card elevation='4'>
    <v-card-text>
      <v-row>
        <v-col cols='6'>
          <div><span class='font-weight-black'>Title:</span> {{title}}</div>
          <div><span class='font-weight-black'>ID:</span> {{id}}</div>
        </v-col>
        <v-col cols='6'>
          <div><span class='font-weight-black'>Phase:</span> {{phase}}</div>
          <div><span class='font-weight-black'>Hand:</span> {{hand}} of {{hands}}</div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup>
import _get from 'lodash/get'
import { computed, inject } from 'vue'

// inject game and current user
import { gameKey } from '@/composables/keys.js'
const { game, update } = inject(gameKey)

const header = computed(() => _get(game, 'value.header', {}))
const title = computed(() => _get(header, 'value.title', ''))
const id = computed(() => _get(game, 'value.id', ''))
const hand = computed(() => _get(header, 'value.round', 0))
const hands = computed(() => _get(game, 'value.hands', 0))
const phase = computed(() => _get(header, 'value.phase', ''))

</script>
