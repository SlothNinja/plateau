<template>
  <div v-if='teams'>
    Increased bid to "{{exchange}} {{objective}} {{teams}}" for a value of {{bValue}}.
  </div>
  <div v-else>
    Increase bid to "{{exchange}} {{objective}}" for a value of {{bValue}}.
  </div>
</template>

<script setup>
import { gameKey } from '@/composables/keys.js'
import { computed, inject, unref } from 'vue'
import _get from 'lodash/get'
import { bidValue } from '@/composables/bid.js'
const props = defineProps(['message'])

const game = inject(gameKey)
const bid = computed(() => _get(props, 'message.Data.Bid', {}))
const exchange = computed(() => _get(unref(bid), 'Exchange', ''))
const objective = computed(() => _get(unref(bid), 'Objective', ''))
const teams = computed(() => _get(unref(bid), '.Teams', ''))
const bValue = computed(() => bidValue(game, bid))

</script>
