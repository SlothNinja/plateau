<template>
  <v-card height='250'>
    <v-card-title>
      {{title}}
    </v-card-title>
    <v-card-text>
      <CardDisplay :height='height' v-model:cards='cards' :dTeam='dTeam' />
    </v-card-text>
  </v-card>
</template>

<script setup>
import CardDisplay from '@/components/Game/CardDisplay.vue'
import { gameKey } from '@/composables/keys.js'
import { computed, inject, unref } from 'vue'
import _get from 'lodash/get'

const props = defineProps(['height'])

const game = inject(gameKey)
const header = computed(() => _get(unref(game), 'Header', {}))
const state = computed(() => _get(unref(game), 'State', {}))
const cards = computed(() => _get(unref(state), `Tricks[${unref(header).Turn}].Cards`, []))
const trick = computed(() => (_get(unref(header), 'Turn', 0) + 1))
const dTeam = computed(() => _get(unref(state), 'DeclarersTeam', []))

const title = computed(() => {
  if (unref(trick) == 14) {
    return 'Talon:'
  }
  return `Trick: ${unref(trick)}`
})

</script>
