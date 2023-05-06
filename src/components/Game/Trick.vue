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
import { computed, inject, unref } from 'vue'
import _get from 'lodash/get'

const game = inject(gameKey)
const props = defineProps(['height'])
const cards = computed(() => _get(unref(game), `Tricks[${unref(game).Turn}].Cards`, []))
const trick = computed(() => (_get(unref(game), 'Turn', 0) + 1))

const title = computed(() => {
  if (unref(trick) == 14) {
    return 'Talon:'
  }
  return `Trick: ${unref(trick)}`
})

</script>
