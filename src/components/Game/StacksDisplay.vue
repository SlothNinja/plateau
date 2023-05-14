<template>
  <v-card elevation='4'>
    <v-card-title>{{title}}</v-card-title>
    <v-card-text class='h-100 w-100' >
      <Stacks :stacks='stacks' :height='height' />
    </v-card-text>
  </v-card>
</template>

<script setup>
// components
import Stacks from '@/components/Game/Stacks'

// lodash
import _sortBy from 'lodash/sortBy'
import _size from 'lodash/size'
import _map from 'lodash/map'
import _takeRight from 'lodash/takeRight'
import _includes from 'lodash/includes'
import _remove from 'lodash/remove'
import _get from 'lodash/get'

// vue
import { computed, ref, inject, unref, watch } from 'vue'

// composables
import { useCardValue } from '@/composables/cardValue.js'
import { useIsCP, usePlayerByUser } from '@/composables/player.js'
import { useUserByIndex } from '@/composables/user.js'
import { cuKey, gameKey } from '@/composables/keys.js'

const props = defineProps({
  height: [ Number, String],
  title: String,
  multi: Number,
  selected: Array,
  sort: Boolean,
  stacks: Array,
})
const emit = defineEmits(['update:cards', 'update:selected'])

const player = computed(() => usePlayerByUser(game, cu))
const hand = computed({
  get() {
    return _get(props, 'cards', [])
  },
  set(value) {
    emit('update:cards', value)
  }
})


function nameFor(card) {
  const pid = _get(card, 'PlayedBy', 0)
  if (pid <= 0) {
    return ''
  }
  const user = useUserByIndex(game, pid - 1)
  return _get(user, 'Name', '')
}

const hover = ref([])

const cu = inject(cuKey)
const game = inject(gameKey)


const isCP = computed(() => useIsCP(game, cu))

function hovered(index, state) {
  const lastIndex = _size(unref(sorted)) - 1
  if (index != lastIndex && state == true) {
    hover.value[index] = true
    hover.value[lastIndex] = false
    return
  }
  if (index != lastIndex && state == false) {
    hover.value[index] = false
  }
  hover.value[lastIndex] = true
}

function initHover() {
  const lastIndex = _size(unref(sorted)) - 1
  hover.value = _map(unref(sorted), () => false)
  hover.value[lastIndex] = true
}

const hoverstyle = 'overflow:visible'
const nohoverstyle = 'overflow:hidden'

const sorted = computed(() => {
  if (props.sort) {
    return _sortBy(hand.value, [ card => card.Suit, useCardValue ])
  }
  return hand.value
})

const handSize = computed(() => _size(unref(sorted)))
const selection = computed({
  get() {
    return _get(props, 'selected', [])
  },
  set(value) {
    emit('update:selected', value)
  }
})

watch( handSize, () => { initHover() } )

const disableSelect = computed(() => {
  return (unref(game).Phase == 'bid') || (!unref(isCP)) || (unref(player).PerformedAction)
})

function select(card) {
  if (unref(disableSelect)) {
    return
  }

  if (_includes(unref(selection), card)) {
    _remove(selection.value, card)
  } else {
    selection.value.push(card)
    selection.value = _takeRight(selection.value, props.multi)
  }
}

function isSelected(card) {
  if (_includes(unref(selection), card)) {
    return 'selected'
  }
}

</script>

<style lang='sass'>
.card
  min-width:80px

.selected
  margin-top:-5%

.playable
  margin-top:-2.5%
</style>
