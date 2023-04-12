<template>
  <v-row>
    <v-col cols='12'>
    </v-col>
  </v-row>
</template>

<script setup>
import Hand from '@/components/Hand.vue'
import { computed, inject } from 'vue'
import { usePlayerFor } from '@/composables/playerFor.js'
import { usePIDFor } from '@/composables/pidFor.js'
import _get from 'lodash/get'

// inject game and current user
import { cuKey, gameKey } from '@/composables/keys.js'
const cu = inject(cuKey)
const game = inject(gameKey)

const pid = computed(() => usePIDFor(game, cu))
const player = computed(() => usePlayerFor(game, cu))
const hand = computed(() => _get(player, 'value.hand', []))
const height = 170
</script>
