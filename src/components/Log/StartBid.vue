<template>
  <PlayerEntry class='ma-1' :pid='pid'>
    <template #toolbar-title>
      Hand: {{handNumber}}
    </template>
    Starts bid and leads first trick.
  </PlayerEntry>
</template>

<script setup>
import PlayerEntry from '@/snvue/components/Log/PlayerEntry'
import { gameKey } from '@/snvue/composables/keys'
import { computed, inject, unref } from 'vue'
import _get from 'lodash/get'

const props = defineProps(['data'])

const game = inject(gameKey)

const pid = computed(() => _get(unref(props), 'data.PID', ''))
const handNumber = computed(() => _get(unref(props), 'data.HandNumber', -1))
</script>
