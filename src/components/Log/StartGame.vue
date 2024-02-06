<template>
  <GameEntry class='ma-1'>

    <div>
      Good luck.  Have fun.
    </div>
    <div class="mt-2 d-flex justify-space-around">
      <UserButton
          v-for='(user, index) in users'
          :key='index'
          :user='user'
          :size='36'
          variant='bottom'
          />
    </div>

  </GameEntry>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton'
import GameEntry from '@/snvue/components/Log/GameEntry'

import { gameKey } from '@/snvue/composables/keys'
import { useUserByIndex } from '@/snvue/composables/user'

import { computed, inject, unref } from 'vue'

import _get from 'lodash/get'
import _map from 'lodash/map'

const props = defineProps(['data'])

const game = inject(gameKey)

const pids = computed(() => _get(unref(props), 'data.PIDS', []))
const header = computed(() => _get(unref(game), 'Header', {}))
const users = computed(() => _map(unref(pids), pid => useUserByIndex(header, pid-1)))
</script>
