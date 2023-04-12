<template>
  <v-card>
    <v-card-title>
      {{title}}
    </v-card-title>
    <v-card-text>
      <CardDisplay :height='height' v-model:cards='cards' />
    </v-card-text>
  </v-card>
</template>

<script setup>
import CardDisplay from '@/components/Game/CardDisplay.vue'
import { gameKey } from '@/composables/keys.js'
import { computed, inject } from 'vue'
import _get from 'lodash/get'

const { game, updateGame } = inject(gameKey)
const props = defineProps(['height'])
const cards = computed(() => _get(game, `value.state.tricks[${game.value.header.turn}].cards`, []))
const trick = computed(() => (_get(game, 'value.header.turn', 0) + 1))

const title = computed(() => {
  if (trick.value == 14) {
    return 'Talon:'
  }
  return `Trick: ${trick.value}`
})

</script>
