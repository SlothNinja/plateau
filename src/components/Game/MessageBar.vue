<template>
  <v-card elevation='4' class='h-100 d-flex justify-center align-center'>
    <div>{{message}}</div>
  </v-card>
</template>

<script setup>
// composables
import { cuKey, gameKey } from '@/snvue/composables/keys.js'
import { useCP, useCPID, useIsCP, usePIDForUser } from '@/composables/player.js'
import { useUserByIndex } from '@/snvue/composables/user.js'

// vue
import { computed, inject, unref } from 'vue'

// lodash
import _get from 'lodash/get'

const cu = inject(cuKey)
const game = inject(gameKey)

const header = computed(() => _get(unref(game), 'Header', {}))
const phase = computed(() => _get(unref(header), 'Phase', ''))
const cp = computed(() => useCP(game))
const cpid = computed(() => useCPID(header))
const isCP = computed(() => useIsCP(header, cu))
const uIndex = computed(() => (unref(cpid) - 1))
const user = computed(() => useUserByIndex(header, uIndex))
const waitMessage = computed(() => (`Please wait for ${unref(user).Name} to take a turn.`))
const numPlayers = computed(() => _get(unref(header), 'NumPlayers', 0))

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

      if (unref(numPlayers) == 2) {
        return 'Select two cards from your hand to return to talon.'
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

function declarer(player) {
  return _first(_get(unref(game), 'DeclarersTeam', [])) == unref(player).ID
}

</script>
