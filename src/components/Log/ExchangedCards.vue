<template>
  <PlayerEntry class='ma-1' :pid='pid'>
  
    <template #toolbar-title>
      Hand: {{handNumber}}
    </template>

    <div>
      Exchanged cards with talon.
    </div>
  </PlayerEntry>
</template>

<script setup>
import PlayerEntry from '@/snvue/components/Log/PlayerEntry'
import { computed, unref } from 'vue'
import _get from 'lodash/get'

const props = defineProps(['data'])
const pid = computed(() => _get(unref(props), 'data.PID', ''))
const handNumber = computed(() => _get(unref(props), 'data.HandNumber', -1))
</script>
