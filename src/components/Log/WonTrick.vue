<template>
  <div class='mt-2 d-flex align-center'>
    <div>
      <UserButton :user='user' :size='36' :color='useColorFor(entry.DeclarersTeam, pid)' />
    </div>
    <div class='w-100 ml-1'>
      won trick.
    </div>
  </div>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton'
import { useColorFor } from '@/composables/color'

import { gameKey } from '@/composables/keys'
import { useUserByPID } from '@/composables/user'

import { computed, inject, unref } from 'vue'

import _get from 'lodash/get'

const props = defineProps(['message', 'entry'])

const game = inject(gameKey)

const header = computed(() => _get(unref(game), 'Header', {}))
const pid = computed(() => _get(props, 'message.Data.PID', 0))
const user = computed(() => useUserByPID(unref(header), unref(pid)))
</script>
