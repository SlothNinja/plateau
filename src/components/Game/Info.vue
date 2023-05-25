<template>
  <v-card elevation='4' class='h-100'>
    <v-card-text>
      <v-row>
        <v-col cols='6'>
          <div><span class='font-weight-black'>Title:</span> {{title}}</div>
          <div><span class='font-weight-black'>ID:</span> {{id}}</div>
        </v-col>
        <v-col cols='6'>
          <div><span class='font-weight-black'>Phase:</span> {{phase}}</div>
          <div><span class='font-weight-black'>Hand:</span> {{hand}} of {{hands}}</div>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup>
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import { computed, inject, unref } from 'vue'
import { useHands } from '@/composables/hands'
import { useRoute } from 'vue-router'

const route = useRoute()

// inject game and current user
import { gameKey } from '@/composables/keys'
const game = inject(gameKey)

const title = computed(() => _get(unref(game), 'Title', ''))
const id = computed(() => _get(unref(route), 'params.id', ''))
const hand = computed(() => _get(unref(game), 'Round', 0))
const hands = computed(() => (useHands(game)))
const phase = computed(() => _get(unref(game), 'Phase', ''))

</script>
