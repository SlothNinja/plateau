<template>
  <div class='d-flex justify-center align-center ma-2' style='height:2em'>
    <v-btn v-if='canSubmit' @click='submit' size='small' color='green'>Submit</v-btn>
  </div>
  <CardDisplay sort v-bind='$attrs' :height='height' :multi='multi' v-model:cards='hand' v-model:selected='selected' />
</template>

<script setup>
// components
import CardDisplay from '@/components/Game/CardDisplay.vue'

// lodash
import _size from 'lodash/size'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

// vue
import { computed, ref, inject, unref, watch } from 'vue'
import { useIsCP, usePlayerByUser } from '@/composables/player.js'
import { usePut } from '@/composables/fetch.js'
import { useRoute } from 'vue-router'

const route = useRoute()

// composables
import { cuKey, gameKey, snackKey } from '@/composables/keys.js'

const player = computed(() => usePlayerByUser(game, cu))

const hand = computed({
  get() {
    return _get(unref(player), 'Hand', [])
  },
  set(value) {
    player.value.hand = value
  }
})

const props = defineProps([ 'height' ])
const hover = ref([])

const cu = inject(cuKey)
const game = inject(gameKey)

const isCP = computed(() => useIsCP(game, cu))

const phase = computed(() => _get(unref(game), 'Phase', ''))
const performedAction  = computed(() => _get(unref(player), 'PerformedAction', false))

const selected = ref([])

watch(
  selected,
  () => {
    if ((unref(multi) == 1) && unref(canSubmit)) {
      submit()
    }
  }
)

const canSubmit = computed(() => (
  unref(isCP) &&
  !unref(performedAction) &&
  ((unref(phase) == 'card exchange') || (unref(phase) == 'card play')) &&
  _size(unref(selected)) == unref(multi)
))

const multi = computed(() => ( unref(phase) == 'card exchange' ? 3 : 1 ))

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

/////////////////////////////////////
// Submit bid to server
function submit() {
  const action = unref(phase) == 'card exchange' ? 'exchange': 'play'
  const { response, error } = usePut(`/sn/game/${action}/${route.params.id}`, unref(selected))
  selected.value = []

  watch(response, () => update(response))
}

function update(response) {
    // const g = _get(response, 'value.game', {})
    // if (!_isEmpty(g)) {
    //   updateGame(g)
    // }
    const msg = _get(unref(response), 'Message', '')
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
