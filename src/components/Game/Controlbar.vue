<template>
  <v-row no-gutters>
    <v-col align="center">

      <!-- Reset control -->
      <v-tooltip location='bottom' text='Reset' color="info" :disabled='!canReset'>
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon='mdi-close' :disabled='!canReset' @click="action('reset')" />
        </template>
      </v-tooltip>

      <!-- Rollback control -->
      <v-tooltip location='bottom' text='Rollback' color="info" :disabled='!canRollback' v-if='admin' >
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon='mdi-step-backward' :disabled='!canRollback' @click="action('rollback', { rev: stack.Committed })" />
        </template>
      </v-tooltip>

      <!-- Rollforward control -->
      <v-tooltip location='bottom' text='Rollforward' color="info" :disabled='!canRollforward' v-if='admin'>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon='mdi-step-forward' @click="action('rollforward', { rev: stack.Committed })" />
        </template>
      </v-tooltip>

      <!-- Undo control -->
      <v-tooltip location='bottom' text='Undo' color="info" :disabled='!canUndo'>
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon='mdi-undo' :disabled='!canUndo' @click="action('undo')" />
        </template>
      </v-tooltip>

      <!-- Redo control -->
      <v-tooltip location='bottom' text='Redo' color="info" :disabled='!canRedo'>
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon='mdi-redo' :disabled='!canRedo' @click="action('redo')" />
        </template>
      </v-tooltip>

      <!-- Finish control -->
      <v-tooltip location='bottom' text='Finish' color="info" :disabled='!canFinish' >
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon='mdi-check' :disabled='!canFinish' @click="action(finishPath)" />
        </template>
      </v-tooltip>

    </v-col>
  </v-row>
</template>

<script setup>
import { computed, inject, watch, unref } from 'vue'

// lodash
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

// useRoute
import { useRoute } from 'vue-router'
const route = useRoute()

// inject current user
import { cuKey, stackKey } from '@/composables/keys.js'
const cu = inject(cuKey)
const admin = _get(unref(cu), 'Admin', false)

//////////////////////////////////////
// Snackbar
import { snackKey } from '@/composables/keys.js'
const { snackbar, updateSnackbar } = inject(snackKey)
const stack = inject(stackKey)

// Inject game
import { gameKey } from '@/composables/keys'
const game = inject(gameKey)
const header = computed(() => _get(unref(game), 'Header', {}))

import { useCP, useIsCP } from '@/composables/player.js'
const isCP = computed(() => useIsCP(header, cu))
const performedAction = computed(() => _get(unref(useCP(game)), 'PerformedAction', false))
const running = computed(() => (_get(unref(game), 'Header.Status', '') == 'running'))

const canFinish = computed(() => (unref(running) && unref(isCP) && unref(performedAction)))
const canUndo = computed(() => (unref(running) && unref(isCP) && (unref(stack).Current > unref(stack).Committed)))
const canRedo = computed(() => (unref(running) && unref(isCP) && (unref(stack).Current < unref(stack).Updated)))
const canReset = computed(() => (unref(running) && unref(isCP)))

const canRollback = computed(() => (unref(admin) && (unref(stack).Current == unref(stack).Committed) && (unref(stack).Committed) > 0))
const canRollforward = computed(() => (unref(admin) && (unref(stack).Current == unref(stack).Committed)))

import { usePut } from '@/composables/fetch.js'
function action(path, data) {
  let url = `/sn/game/${path}/${route.params.id}`
  if (process.env.NODE_ENV == 'development') {
    const backend = import.meta.env.VITE_PLATEAU_BACKEND
    url = `${backend}sn/game/${path}/${route.params.id}`
  }
  const { state, isReady, isLoading } = usePut(url, data)

  watch( state, () => {
    const msg = _get(unref(state), 'Message', '')
    if (!_isEmpty(msg)) {
      updateSnackbar(msg)
    }
  })
}

const finishPath = computed(
  ()=> {
    switch(_get(unref(game), 'Header.Phase', '')) {
      case 'bid':
        return 'finish/bid'
      case 'pass':
        return 'finish/pass'
      case 'card exchange':
        return 'finish/exchange'
      case 'increase objective':
        return 'finish/objective'
      case 'card play':
        return 'finish/play'
      default:
        return 'finish'
    }
  }
)
</script>
