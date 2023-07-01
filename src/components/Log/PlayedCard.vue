<template>
  <CardDisplay :height='height' v-model:cards='cards' />
</template>

<script setup>
import CardDisplay from '@/components/Game/CardDisplay.vue'
import { gameKey } from '@/composables/keys.js'
import { computed, inject, unref } from 'vue'
import _get from 'lodash/get'
import _each from 'lodash/each'

const props = defineProps(['message'])
const game = inject(gameKey)
const height = '120'
const cards = computed(
  () => {
    let collect = []
    _each(props.message.Data, (v, k) => collect[k] = v)
    return collect
  }
)

</script>
