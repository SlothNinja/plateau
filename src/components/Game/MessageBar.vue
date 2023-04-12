<template>
  <v-card elevation='4' class='d-flex justify-center align-center'>
    <div>{{message}}</div>
  </v-card>
</template>

<script setup>
// composables
import { cuKey, gameKey } from '@/composables/keys.js'
import { useCP, useCPID, useIsCP, usePIDForUser } from '@/composables/player.js'
import { useUserByIndex } from '@/composables/user.js'

// vue
import { computed, inject } from 'vue'

// lodash
import _get from 'lodash/get'

const cu = inject(cuKey)
const { game, update } = inject(gameKey)

const header = computed(() => _get(game, 'value.header', {}))
const phase = computed(() => _get(header, 'value.phase', {}))
const cp = computed(() => useCP(game))
const cpid = computed(() => useCPID(game))
const isCP = computed(() => useIsCP(game, cu))
const uIndex = computed(() => (cpid.value - 1))
const user = computed(() => useUserByIndex(header, uIndex))
const waitMessage = computed(() => (`Please wait for ${user.value.name} to take a turn.`))

const message = computed(() => {
  switch(phase.value) {
    case 'bid':
      if (!isCP.value) {
        return waitMessage.value
      }

      if (!cp.value.bid) {
        return 'Submit bid or pass.'
      }

      if (cp.value.performedAction) {
        return 'Finish turn by selecting above check mark.'
      }
      return ''
    case 'card exchange':
      if (!isCP.value) {
        return waitMessage.value
      }

      if (cp.value.performedAction) {
        return 'Finish turn by selecting above check mark.'
      }

      return 'Select three cards from your hand to return to talon.'
    case 'increase objective':
      if (!isCP.value) {
        return waitMessage.value
      }

      if (cp.value.performedAction) {
        return 'Finish turn by selecting above check mark.'
      }

      return 'You may increase the objective of the bid.'
    case 'card play':
      if (!isCP.value) {
        return waitMessage.value
      }

      if (cp.value.performedAction) {
        return 'Finish turn by selecting above check mark.'
      }

      return 'Play card to trick'
    default:
      return ''
  }
})

</script>
