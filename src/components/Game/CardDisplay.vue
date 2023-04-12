<template>
  <v-sheet :height='height'>
    <div class='d-flex justify-center align-center h-100 w-100' >
      <div
          v-for='(card, index) in sorted'
          :key='index'
          @mouseover='hovered(index, true)'
          @mouseleave='hovered(index, false)'
          :style='hover[index] ? hoverstyle : nohoverstyle'
          :class='isSelected(card)'
          >
          <Card 
          @click='select(card)'
          :rank='card.rank'
          :suit='card.suit'
          :width='height / 2.0'
          :text='nameFor(card)'
          />
      </div>
    </div>
  </v-sheet>
</template>

<script setup>
// components
import Card from '@/components/Game/Card.vue'

// lodash
import _sortBy from 'lodash/sortBy'
import _size from 'lodash/size'
import _map from 'lodash/map'
import _takeRight from 'lodash/takeRight'
import _includes from 'lodash/includes'
import _remove from 'lodash/remove'
import _get from 'lodash/get'

// vue
import { computed, ref, inject, watch } from 'vue'

// composables
import { useCardValue } from '@/composables/cardValue.js'
import { useIsCP, usePlayerByUser } from '@/composables/player.js'
import { useUserByIndex } from '@/composables/user.js'
import { cuKey, gameKey } from '@/composables/keys.js'

const props = defineProps({
  height: [ Number, String],
  multi: Number,
  cards: Array,
  selected: Array,
  sort: Boolean,
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
  const pid = _get(card, 'playedBy', 0)
  if (pid <= 0) {
    return ''
  }
  const user = useUserByIndex(_get(game, 'value.header', {}), pid - 1)
  return _get(user, 'name', '')
}

const hover = ref([])

const cu = inject(cuKey)
const { game, updateGame } = inject(gameKey)


const isCP = computed(() => useIsCP(game, cu))

function hovered(index, state) {
  const lastIndex = _size(sorted.value) - 1
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
  const lastIndex = _size(sorted.value) - 1
  hover.value = _map(sorted.value, () => false)
  hover.value[lastIndex] = true
}

const hoverstyle = 'overflow:visible'
const nohoverstyle = 'overflow:hidden'

const sorted = computed(() => {
  if (props.sort) {
    return _sortBy(hand.value, [ card => card.suit, useCardValue ])
  }
  return hand.value
})

const handSize = computed(() => _size(sorted.value))
const selection = computed({
  get() {
    return _get(props, 'selected', [])
  },
  set(value) {
    emit('update:selected', value)
  }
})

watch( handSize, () => { initHover() } )

function select(card) {
  if (!isCP.value || player.value.performedAction) {
    return
  }

  if (_includes(selection.value, card)) {
    _remove(selection.value, card)
  } else {
    selection.value.push(card)
    selection.value = _takeRight(selection.value, props.multi)
  }
}

function isSelected(card) {
  if (_includes(selection.value, card)) {
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
