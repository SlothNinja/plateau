<template>
  <div class='mt-2 d-flex align-center'>
    <div>
      <UserButton :user='user' :size='36' />
    </div>
    <div class='w-100 ml-1'>
      won trick.
    </div>
  </div>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton'

import { gameKey } from '@/composables/keys'
import { useUserByIndex } from '@/composables/user'

import { computed, inject, unref } from 'vue'

import _get from 'lodash/get'

const props = defineProps(['message'])

const game = inject(gameKey)

const pid = computed(() => _get(props, 'message.Data.PID', 0))
const user = computed(() => useUserByIndex(unref(game), unref(pid)-1))
</script>
