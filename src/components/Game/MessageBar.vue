<template>
  <v-card elevation='4' class='h-100 d-flex justify-center align-center'>
    <div>{{message}}</div>
  </v-card>
</template>

<script setup>
// composables
import { cuKey, gameKey } from '@/composables/keys.js'
import { useCP, useCPID, useIsCP, usePIDForUser } from '@/composables/player.js'
import { useUserByIndex } from '@/composables/user.js'

// vue
import { computed, inject, unref } from 'vue'

// lodash
import _get from 'lodash/get'

const cu = inject(cuKey)
const game = inject(gameKey)

const phase = computed(() => _get(unref(game), 'Phase', ''))
const cp = computed(() => useCP(game))
const cpid = computed(() => useCPID(game))
const isCP = computed(() => useIsCP(game, cu))
const uIndex = computed(() => (unref(cpid) - 1))
const user = computed(() => useUserByIndex(game, uIndex))
const waitMessage = computed(() => (`Please wait for ${unref(user).Name} to take a turn.`))

const message = computed(() => {
  switch(unref(phase)) {
    case 'bid':
      if (!unref(isCP)) {
        return unref(waitMessage)
      }

      if (unref(cp).PerformedAction) {
        return 'Finish turn by selecting above check mark.'
      }
      return 'Submit bid or pass.'
    case 'card exchange':
      if (!unref(isCP)) {
        return unref(waitMessage)
      }

      if (unref(cp).PerformedAction) {
        return 'Finish turn by selecting above check mark.'
      }

      return 'Select three cards from your hand to return to talon.'
    case 'pick partner':
      if (!unref(isCP)) {
        return unref(waitMessage)
      }

      if (unref(cp).PerformedAction) {
        return 'Finish turn by selecting above check mark.'
      }

      return 'Select card to select partner.'
    case 'increase objective':
      if (!unref(isCP)) {
        return unref(waitMessage)
      }

      if (unref(cp).PerformedAction) {
        return 'Finish turn by selecting above check mark.'
      }

      return 'You may increase the objective of the bid.'
    case 'card play':
      if (!unref(isCP)) {
        return unref(waitMessage)
      }

      if (unref(cp).PerformedAction) {
        return 'Finish turn by selecting above check mark.'
      }

      return 'Play card to trick'
    default:
      return ''
  }
})

</script>
