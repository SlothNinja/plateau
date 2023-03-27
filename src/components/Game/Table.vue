<template>
  <v-sheet elevation='4' rounded color='green' class='d-flex align-center justify-center h-100 w-100'>
    <BidForm v-if="bidPhase" />
    <div v-else>
      Trick play
    </div>
  </v-sheet>
</template>

<script setup>
import BidForm from '@/components/Game/BidForm.vue'
import { computed, inject } from 'vue'
import _get from 'lodash/get'

// inject game 
import { gameKey } from '@/composables/keys.js'
const { game, updateGame } = inject(gameKey)

// const header = computed(() => _get(game, 'value.header', {}))
const phase = computed(() => _get(game, 'value.header.phase', ''))
const bidPhase = computed(() => (phase.value == 'bid'))

</script>
