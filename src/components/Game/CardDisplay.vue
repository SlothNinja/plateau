<template>
  <div class='d-flex justify-center align-center h-100 w-100' >
    <div
        v-for='(card, index) in sorted'
        :key='index'
        @mouseover='hovered(index, true)'
        @mouseleave='hovered(index, false)'
        :style='hover[index] ? hoverstyle : nohoverstyle'
        :class='isSelected(card)'
        >

        <!-- Display Card Stacks -->
      <div v-if='_isArray(card)' :style='`height:${height*1.1}px`'>
        <Card 
           class='mx-1'
           v-if='_size(card) > 1'
           :rank='_first(card).Rank'
           :suit='_first(card).Suit'
           :width='height / 2.0'
           :text='nameFor(_first(card))'
           :textcolor='useColorFor(dTeam, _first(card).PlayedBy)'
           />
        <Card 
           class='mx-1'
           v-if='_size(card) > 1'
           @click='select(_last(card))'
           :rank='_last(card).Rank'
           :suit='_last(card).Suit'
           :width='height / 2.0'
           :text='nameFor(_last(card))'
           :textcolor='useColorFor(dTeam, _last(card).PlayedBy)'
           stacked
           />
        <Card 
           class='mx-1'
           v-if='_size(card) == 1'
           @click='select(_first(card))'
           :rank='_first(card).Rank'
           :suit='_first(card).Suit'
           :width='height / 2.0'
           :text='nameFor(_first(card))'
           :textcolor='useColorFor(dTeam, _first(card).PlayedBy)'
           />
      </div>

      <!-- Display Cards in Hand -->
      <div v-else>
        <Card 
           class='mx-1'
           @click='select(card)'
           :rank='card.Rank'
           :suit='card.Suit'
           :width='height / 2.0'
           :text='nameFor(card)'
           :textcolor='useColorFor(dTeam, card.PlayedBy)'
           />
      </div>
    </div>
  </div>
</template>

<script setup>
// components
import Card from '@/components/Game/Card.vue'

// lodash
import _sortBy from 'lodash/sortBy'
import _size from 'lodash/size'
import _last from 'lodash/last'
import _first from 'lodash/first'
import _map from 'lodash/map'
import _takeRight from 'lodash/takeRight'
import _includes from 'lodash/includes'
import _remove from 'lodash/remove'
import _isArray from 'lodash/isArray'
import _get from 'lodash/get'

// vue
import { computed, ref, inject, unref, watch } from 'vue'

// composables
import { useCardValue } from '@/composables/cardValue'
import { useIsCP, usePlayerByUser } from '@/composables/player'
import { useUserByPID } from '@/composables/user'
import { cuKey, gameKey } from '@/composables/keys'
import { useColorFor } from '@/composables/color'

const cu = inject(cuKey)
const game = inject(gameKey)

const props = defineProps({
  height: [ Number, String],
  multi: Number,
  cards: Array,
  selected: Array,
  sort: Boolean,
  dTeam: Array,
})

const emit = defineEmits(['update:selected'])

const player = computed(() => usePlayerByUser(game, cu))
const header = computed(() => _get(unref(game), 'Header', {}))

function nameFor(card) {
  const pid = _get(card, 'PlayedBy', 0)
  if (pid <= 0) {
    return ''
  }
  const user = useUserByPID(header, pid)
  return _get(user, 'Name', '')
}

const hover = ref([])


const isCP = computed(() => useIsCP(header, cu))

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
    return _sortBy(unref(props.cards), [ card => card.Suit, useCardValue ])
  }
  return unref(props.cards)
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
  return (unref(header).Phase == 'bid') || (!unref(isCP)) || (unref(player).PerformedAction)
})

function select(card) {
  if (unref(disableSelect)) {
    return
  }

  if (_includes(unref(selection), card)) {
    _remove(selection.value, card)
  } else if ((unref(header).Phase != 'card exchange') ||
    ((unref(header).Phase != 'card exchange') && (_includes(unref(player).Hand, card)))) {
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
