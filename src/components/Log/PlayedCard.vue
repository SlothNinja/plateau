<template>
  <CardDisplay :height='height' v-model:cards='cards' :dTeam='entry.DeclarersTeam' />
</template>

<script setup>
import CardDisplay from '@/components/Game/CardDisplay'
import { computed, unref } from 'vue'
import _each from 'lodash/each'

const props = defineProps({
  message: { type: Object, required: true },
  entry: { type: Object, default: {} },
})

const height = '120'
const cards = computed(
  () => {
    let collect = []
    _each(props.message.Data, (v, k) => collect[k] = v)
    return collect
  }
)

</script>
