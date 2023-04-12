<template>
  <v-container>
    <v-row dense>
      <v-col cols='6'>
        <v-row dense>
          <v-col cols='12'>
            <Info />
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <MessageBar />
          </v-col>
        </v-row>
        <v-row dense v-if='showTrick'>
          <v-col cols='12'>
            <Trick height='170' />
          </v-col>
        </v-row>
        <v-row dense v-if='showForm'>
          <v-col cols='12'>
            <BidForm />
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <Hand class='pa-2' height='170' />
          </v-col>
        </v-row>
      </v-col>
      <v-col cols='6'>
        <v-row dense>
          <v-col cols='12'>
            <Bids />
          </v-col>
        </v-row>
        <v-row dense>
          <v-col cols='12'>
            <Board :tricks='tricks' :declarersTeam='declarersTeam' />
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
// components
import Bids from '@/components/Game/Bids.vue'
import BidForm from '@/components/Game/BidForm.vue'
import Trick from '@/components/Game/Trick.vue'
import Hand from '@/components/Game/Hand.vue'
import Info from '@/components/Game/Info.vue'
import Board from '@/components/Game/Board.vue'
import Table from '@/components/Game/Table.vue'
import MessageBar from '@/components/Game/MessageBar.vue'

// composables
import { cuKey, gameKey } from '@/composables/keys.js'
import { usePlayerByUser, useIsCP } from '@/composables/player.js'
import { useFetch } from '@/composables/fetch.js'

// vue
import { useRoute } from 'vue-router'
import { computed, provide, inject, ref, watch, onMounted, onUnmounted } from 'vue'

// lodash
import _get from 'lodash/get'

const cu = inject(cuKey)

const player = computed(() => usePlayerByUser(game, cu))
const hand = computed(() => _get(player, 'value.hand', []))

const height = 170

const route = useRoute()


const { game, updateGame } = inject(gameKey)

// const declarersCards = computed(() => _get(game, 'value.state.declarersCards', []))
// const opposersCards = computed(() => _get(game, 'value.state.opposersCards', []))
const tricks = computed(() => _get(game, 'value.state.tricks', []))
const declarersTeam = computed(() => _get(game, 'value.state.declarersTeam', []))

function fetch() {
  const { data, error } = useFetch(`/sn/game/show/${route.params.id}`)

  watch(data, () => {
    updateGame(_get(data, 'value.game', {}))
  })
}

const isCP = computed(() => (useIsCP(game, cu)))

const phase = computed(() => _get(game, 'value.header.phase', ''))
const showForm = computed(() => {
  return isCP.value && (phase.value == 'bid' || phase.value == 'increase objective')
})

const showTrick = computed(() => {
  return phase.value == 'card play'
})

onMounted(() => {
  window.addEventListener('blur', () => (shown.value = false))
  window.addEventListener('focus', () => (shown.value = true))
  fetch()
})

onUnmounted(() => {
  window.removeEventListener('blur', () => (shown.value = false))
  window.removeEventListener('focus', () => (shown.value = true))
})

const shown = ref(false)
watch(shown, fetch)

</script>
