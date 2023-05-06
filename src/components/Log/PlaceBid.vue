<template>
  <li class='ml-6'>
    Placed a bid of "{{exchange}} {{objective}} {{teams}}" for a value of {{bValue}}.
  </li>
</template>

<script setup>
import { gameKey } from '@/composables/keys.js'
import { computed, inject } from 'vue'
import _get from 'lodash/get'
import { bidValue } from '@/composables/bid.js'
const props = defineProps(['message'])

const game = inject(gameKey)
const exchange = computed(() => _get(props.message.bid, 'exchange', ''))
const objective = computed(() => _get(props.message.bid, 'objective', ''))
const teams = computed(() => _get(props.message.bid, 'teams', ''))
const bValue = computed(() => bidValue(game, props.message.bid))

</script>
