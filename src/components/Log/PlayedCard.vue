<template>
  <GameEntry class='ma-1'>

    <template #toolbar-title>
      Hand: {{handNumber}} Trick: {{trickNumber}}
    </template>

    <CardDisplay :height='height' v-model:cards='data.Trick.Cards' :dTeam='data.DeclarersTeam' />
    <div v-if='showWinner' class='mt-2 d-flex align-center'>
      <div>
        <UserButton :user='user' :size='36' :color='useColorFor(data.DeclarersTeam, pid)' />
      </div>
      <div class='w-100 ml-1'>
        won trick.
      </div>
    </div>

  </GameEntry>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton'
import CardDisplay from '@/components/Game/CardDisplay'
import GameEntry from '@/snvue/components/Log/GameEntry'

import { computed, inject, unref, watchEffect, onMounted } from 'vue'
import { useColorFor } from '@/composables/color'
import { gameKey } from '@/snvue/composables/keys'
import { useUserByPID } from '@/snvue/composables/user'
import _get from 'lodash/get'

const height = '120'

const props = defineProps(['data'])
const model = defineModel()

const game = inject(gameKey)
const header = computed(() => _get(unref(game), 'Header', {}))
const pid = computed(() => _get(props, 'data.WonTrick', -1))
const user = computed(() => useUserByPID(unref(header), unref(pid)))
const showWinner = computed(() => unref(pid) != -1)
const handNumber = computed(() => _get(props, 'data.HandNumber', ''))
const trickNumber = computed(() => _get(props, 'data.TrickNumber', ''))
</script>
