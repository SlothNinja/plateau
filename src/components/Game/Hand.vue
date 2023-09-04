<template>
  <v-card elevation='4'>
    <v-card-title class='d-flex'>
      <div>{{useNameFor(game, pid)}}</div>
      <div class='text-center w-100'>
        <v-btn v-if='canSubmit' @click='submit' size='small' color='green'>Submit</v-btn>
      </div>
    </v-card-title>
    <v-card-text class='h-100 w-100'>
      <Stacks
          v-if='showStack'
          :height='height'
          :stacks='myStacks'
          v-model:selected='selected'
          />
      <CardDisplay
          v-if='pickPartner'
          v-bind='$attrs'
          :height='height'
          :multi='multi'
          :cards='pickPartnerCards'
          v-model:selected='selected'
          />
      <CardDisplay
          sort
          v-bind='$attrs'
          :height='height'
          :multi='multi'
          :cards='hand'
          v-model:selected='selected'
          />
    </v-card-text>
  </v-card>
</template>

<script setup>
// components
import CardDisplay from '@/components/Game/CardDisplay.vue'
import Stacks from '@/components/Game/Stacks.vue'

// lodash
import _size from 'lodash/size'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import _difference from 'lodash/difference'
import _differenceWith from 'lodash/differenceWith'

// vue
import { computed, ref, inject, unref, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// composables
import { cuKey, gameKey, snackKey } from '@/composables/keys'
import { useIsCP, useCPID, useNameFor, usePlayerByUser } from '@/composables/player'
import { usePut } from '@/composables/fetch'
import { useStackByPID } from '@/composables/stack'

const player = computed(() => usePlayerByUser(game, cu))
const pid = computed(() =>_get(unref(player), 'ID', 1))

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
const myStacks = computed(() => useStackByPID(game, pid))
const numPlayers = computed(() => _get(unref(game), 'NumPlayers', 0))
const showStack = computed(() => unref(numPlayers) == 2)

const phase = computed(() => _get(unref(game), 'Phase', ''))
const performedAction  = computed(() => _get(unref(player), 'PerformedAction', false))

const pickPartner = computed(() => (unref(isCP) && (unref(phase) == 'pick partner')))

const pickPartnerCards = computed(() => _get(unref(game), 'Pick', []))

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
  ((unref(phase) == 'card exchange') || (unref(phase) == 'card play') || (unref(phase) == 'pick partner')) &&
  _size(unref(selected)) == unref(multi)
))

const multi = computed(
  () => {
    switch (unref(phase)) {
      case 'card exchange':
        return unref(numPlayers) == 2 ? 2 : 3
      default:
        return 1
    }
  }
)

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

/////////////////////////////////////
// Submit bid to server
function submit() {
  let action = ''
  switch (unref(phase)) {
    case 'card exchange':
      action = 'exchange'
      break
    case 'pick partner':
      action = 'pick'
      break
    default:
      action = 'play'
  }
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
