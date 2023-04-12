<template>
  <div v-if='canSubmit' class='d-flex justify-center align-center ma-2'>
    <v-btn @click='submit' size='small' color='green'>Submit</v-btn>
  </div>
  <CardDisplay sort v-bind='$attrs' :height='height' :multi='multi' v-model:cards='hand' v-model:selected='game.selected' />
</template>

<script setup>
// components
import CardDisplay from '@/components/Game/CardDisplay.vue'

// lodash
import _size from 'lodash/size'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

// vue
import { computed, ref, inject, watch } from 'vue'
import { useIsCP, usePlayerByUser } from '@/composables/player.js'
import { usePut } from '@/composables/fetch.js'

// composables
import { cuKey, gameKey, snackKey } from '@/composables/keys.js'

const player = computed(() => usePlayerByUser(game, cu))

const hand = computed({
  get() {
    return _get(player, 'value.hand', [])
  },
  set(value) {
    player.value.hand = value
  }
})

const props = defineProps([ 'height' ])
const hover = ref([])

const cu = inject(cuKey)
const { game, updateGame } = inject(gameKey)

const isCP = computed(() => useIsCP(game, cu))

const phase = computed(() => _get(game, 'value.header.phase', ''))
const performedAction  = computed(() => _get(player, 'value.performedAction', false))

const canSubmit = computed(() => (
  isCP &&
  !performedAction.value &&
  ((phase.value == 'card exchange') || (phase.value == 'card play')) &&
  _size(game.value.selected) == multi.value)
)

const multi = computed(() => ( phase.value == 'card exchange' ? 3 : 1 ))

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

/////////////////////////////////////
// Submit bid to server
function submit() {
  const action = phase.value == 'card exchange' ? 'exchange': 'play'
  const { response, error } = usePut(`/sn/game/${action}/${game.value.id}`, game.value.selected)

  watch(response, () => update(response))
}

function update(response) {
    const g = _get(response, 'value.game', {})
    if (!_isEmpty(g)) {
      updateGame(g)
    }
    const msg = _get(response, 'value.message', '')
    if (!_isEmpty(msg)) {
      updateSnackbar(msg, true)
    }
}

</script>

<style lang='sass'>
.card
  min-width:80px

.selected
  margin-top:-5%
</style>
