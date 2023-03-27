<template>
  <v-sheet elevation='4' rounded :height='height'>
    <div class='d-flex justify-center align-center h-100 w-100' >
      <div
          v-for='(card, index) in sorted'
          :key='index'
          @mouseover='hovered(index, true)'
          @mouseleave='hovered(index, false)'
          :style='hover[index] ? hoverstyle : nohoverstyle'
          >
          <Card :rank='card.rank' :suit='card.suit' :width='height / 2.0' />
      </div>
    </div>
  </v-sheet>
</template>

<script setup>
import Card from '@/components/Card.vue'
import _sortBy from 'lodash/sortBy'
import _size from 'lodash/size'
import _map from 'lodash/map'
import { computed, ref, watch } from 'vue'
import { useCardValue } from '@/composables/cardValue.js'

const props = defineProps([ 'hand', 'height' ])
const hover = ref([])

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

const sorted = computed(() => _sortBy(props.hand, [ card => card.suit, useCardValue ]))
const handSize = computed(() => _size(sorted.value))

watch( handSize, () => { initHover() } )

</script>

<style lang='sass'>
.card
  min-width:80px
</style>
