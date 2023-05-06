<template>
  <div>
    Good luck.  Have fun.
    <div class="mt-2 d-flex flex-wrap justify-space-around">
      <UserButton
          v-for='(player, index) in game.state.players'
          :key='index'
          :user='useUserByIndex(game.header, player.id-1)'
          :size='36'
          variant='bottom'
          />
    </div>
  </div>
</template>

<script setup>
import UserButton from '@/components/Common/UserButton.vue'
import { gameKey } from '@/composables/keys.js'
import { inject } from 'vue'
import { useUserByIndex } from '@/composables/user.js'

const game = inject(gameKey)
</script>
