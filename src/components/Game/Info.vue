<template>
  <v-row>
    <v-col cols='6'>
      <v-sheet elevation='4' rounded class='h-100 w-100 pa-3'>
        <div><span class='font-weight-black'>Title:</span> {{title}}</div>
        <div><span class='font-weight-black'>ID:</span> {{id}}</div>
        <div><span class='font-weight-black'>Hand:</span> {{hand}} of {{hands}}</div>
      </v-sheet>
    </v-col>
    <v-col cols='6' align='center'>
      <v-sheet elevation='4' rounded class='h-100 w-100 pa-3'>
        <div class='mb-1'><span class='font-weight-black'>Current Player</span></div>
        <div><UserButton :user='user' :size='32' variant='bottom' /></div>
      </v-sheet>
    </v-col>
  </v-row>
</template>

<script setup>
import UserButton from '@/components/UserButton.vue'
import _get from 'lodash/get'
import { computed, inject } from 'vue'
import { useUser } from '@/composables/user.js'

// inject game and current user
import { cuKey, gameKey } from '@/composables/keys.js'
const cu = inject(cuKey)
const game = inject(gameKey)

const header = computed(() => _get(game, 'value.header', {}))
const title = computed(() => _get(header, 'value.title', ''))
const id = computed(() => _get(game, 'value.id', ''))
const hand = computed(() => _get(header, 'value.round', 0))
const hands = computed(() => _get(game, 'value.hands', 0))
const cpid = computed(() => _get(header, 'value.cpids[0]', -1))
const uIndex = computed(() => (cpid.value - 1))
const user = computed(() => useUser(header, uIndex))

</script>
