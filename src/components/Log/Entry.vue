<template>
  <v-card>
    <v-toolbar height='1em' flat color='green'>
      <v-toolbar-title class='text-subtitle-2'>
        Hand: {{entry.handNumber}}
      </v-toolbar-title>
    </v-toolbar>
    <v-card-text class='d-flex align-center'>
      <div v-if='entry.pid'>
        <UserButton
            :user='useUserByIndex(game.header, entry.pid-1)'
            :size='36'
            />
      </div>
      <div class='w-100'>
        <ul>
        <Message
            v-for='(message, index) in entry.messages'
            :key='index'
            :message='message'
            />
        </ul>
      </div>
    </v-card-text>
    <v-divider></v-divider>
    <div class='text-center text-caption'>
      {{updatedAt}}
    </div>
  </v-card>
</template>

<script setup>
// components
import Message from '@/components/Log/Message'
import UserButton from '@/components/UserButton.vue'

//composables
import { useUserByIndex } from '@/composables/user.js'
import { gameKey } from '@/composables/keys.js'

// vue
import { computed, inject } from 'vue'

// lodash
import _get from 'lodash/get'

const props = defineProps(['entry'])

const { game, updateGame } = inject(gameKey)

const updatedAt = computed(() => {
  var d = _get(props.entry, 'updatedAt', false)
  if (d) {
    return new Date(d).toString()
  }
  return ''
})
</script>
