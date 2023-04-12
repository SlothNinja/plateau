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
      <v-tooltip location='bottom' text='Rollback' color="info" :disabled='!canRollback' v-if='cu.admin' >
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon='mdi-step-backward' :disabled='!canRollback' @click="action('rollback', { rev: game.rev })" />
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

      <!-- Rollforward control -->
      <v-tooltip location='bottom' text='Rollforward' color="info" :disabled='!canRollforward' v-if='cu.admin'>
        <template v-slot:activator="{ props }">
            <v-btn v-bind="props" icon='mdi-step-forward' @click="action('rollforward', { rev: game.rev })" />
        </template>
      </v-tooltip>

      <!-- Finish control -->
      <v-tooltip location='bottom' text='Finish' color="info" :disabled='!canFinish' >
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon='mdi-check' :disabled='!canFinish' @click="action('finish')" />
        </template>
      </v-tooltip>

      <!-- Refresh control -->
      <v-tooltip location='bottom' text='Refresh' color="info">
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" icon='mdi-refresh' @click='refresh' />
        </template>
      </v-tooltip>

    </v-col>
  </v-row>
</template>

<script setup>
import { computed, inject, watch } from 'vue'
import { useFetch, usePut } from '@/composables/fetch.js'
import { useCP, useIsCP } from '@/composables/player.js'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

// inject game and current user
import { cuKey, gameKey, snackKey } from '@/composables/keys.js'
const cu = inject(cuKey)

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

const { game, updateGame } = inject(gameKey)

const running = computed(() => (_get(game, 'value.header.status', '') == 'running'))

const cp = computed(() => useCP(game))
const isCP = computed(() => useIsCP(game, cu))
const performedAction = computed(() => _get(cp, 'value.performedAction', false))

const undoCurrent = computed(() => _get(game, 'value.header.undo.current', 0))
const undoCommitted = computed(() => _get(game, 'value.header.undo.committed', 0))
const undoUpdated = computed(() => _get(game, 'value.header.undo.updated', 0))

const canFinish = computed(() => (running.value && isCP.value && performedAction.value))

const canUndo = computed(() => (running.value && isCP.value && (undoCurrent.value > undoCommitted.value)))
const canRedo = computed(() => (running.value && isCP.value && (undoCurrent.value < undoUpdated.value)))
const canReset = computed(() => (running.value && isCP.value))

const canRollback = computed(() => (cu.value.admin && (undoCurrent.value == undoCommitted.value) && (undoCommitted.value > 0)))
const canRollforward = computed(() => (cu.value.admin && (undoCurrent.value == undoCommitted.value)))

function action(path, data) {
  const { response, error } = usePut(`/sn/game/${path}/${game.value.id}`, data)

  watch( response, () => {
    const g = _get(response, 'value.game', {})
    if (!_isEmpty(g)) {
      updateGame(g)
    }
    const msg = _get(response, 'value.message', '')
    if (!_isEmpty(msg)) {
      updateSnackbar(msg)
    }
  })
}

function refresh() {
  const { data, error } = useFetch(`/sn/game/show/${game.value.id}`)

  watch(data, () => {
    updateGame(_get(data, 'value.game', {}))
  })
}
</script>
