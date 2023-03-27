<template>
  <v-container>
    <v-row>
      <v-col cols='6'>
        <MessageBar />
      </v-col>
      <v-col cols='6'>
        <Info />
      </v-col>
    </v-row>
    <v-row>
      <v-col cols='6'>
        <Table class='h-100 w-100'/>
      </v-col>
      <v-col cols='6'>
        <Board />
        <Hand class='pa-2' height='170' :hand='hand' />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import Hand from '@/components/Hand.vue'
import { usePlayerFor } from '@/composables/playerFor.js'
import { usePIDFor } from '@/composables/pidFor.js'

// inject game and current user
import { cuKey, gameKey } from '@/composables/keys.js'
const cu = inject(cuKey)

const pid = computed(() => usePIDFor(game, cu))
const player = computed(() => usePlayerFor(game, cu))
const hand = computed(() => _get(player, 'value.hand', []))
const height = 170
import Info from '@/components/Game/Info.vue'
import Board from '@/components/Game/Board.vue'
import Table from '@/components/Game/Table.vue'
import PlayerHand from '@/components/PlayerHand.vue'
import MessageBar from '@/components/MessageBar.vue'

// Composables
import { useFetch } from '@/composables/fetch.js'
import { useRoute } from 'vue-router'
import { computed, provide, inject, ref, watch } from 'vue'
import _get from 'lodash/get'

const route = useRoute()

const game = ref({})

const { data, error } = useFetch(`/sn/game/show/${route.params.id}`)

watch(data, () => {
  game.value = _get(data, 'value.game', {})
})

function updateGame(value) {
  game.value = value
}

provide(gameKey, { game, updateGame })

</script>
