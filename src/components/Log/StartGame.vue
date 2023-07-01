<template>
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
</template>

<script setup>
import UserButton from '@/components/Common/UserButton'

import { gameKey } from '@/composables/keys'
import { useUserByIndex } from '@/composables/user'

import { computed, inject, unref } from 'vue'

import _get from 'lodash/get'
import _map from 'lodash/map'

const props = defineProps(['message'])

const game = inject(gameKey)

const pids = computed(() => _get(unref(props), 'message.Data.PIDS', []))
const users = computed(() => _map(unref(pids), pid => useUserByIndex(unref(game), pid-1)))
</script>
