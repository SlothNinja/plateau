<template>
  <v-card>
    <v-toolbar height='32' flat color='green'>
      <v-toolbar-title v-if='showHandNumber' class='text-subtitle-2'>
        <span v-if='showHandNumber'>Hand: {{handNumber}}</span> <span v-if='showTrickNumber'>Trick: {{trickNumber}}</span>
      </v-toolbar-title>
    </v-toolbar>

    <v-card-text v-if='pid' class='d-flex align-center'>
      <div>
        <UserButton :user='useUserByIndex(game, pid-1)' :size='36' />
      </div>
      <div class='w-100'>
        <ul>
          <li class='ml-8' v-for='(message, index) in entry.Lines'>
            <Message :key='index' :message='message' />
          </li>
        </ul>
      </div>
    </v-card-text>

    <v-card-text v-else class='d-flex align-center'>
      <div class='w-100'>
        <div class='ml-8' v-for='(message, index) in entry.Lines'>
          <Message :key='index' :message='message' />
        </div>
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
import UserButton from '@/components/Common/UserButton.vue'

//composables
import { useUserByIndex } from '@/composables/user.js'
import { gameKey } from '@/composables/keys.js'

// vue
import { computed, inject, unref } from 'vue'

// lodash
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import _size from 'lodash/size'
import _first from 'lodash/first'

const props = defineProps(['entry'])

const game = inject(gameKey)

const handNumber = computed(() => _get(props, 'entry.Data.HandNumber', -1))
const showHandNumber = computed(() => (unref(handNumber) > -1))

const trickNumber = computed(() => _get(props, 'entry.Data.TrickNumber', 0))
const showTrickNumber = computed(() => (unref(trickNumber) > 0))

const pid = computed(() => _get(props, 'entry.Data.PID', ''))
const updatedAt = computed(() => {
  var d = _get(props.entry, 'updatedAt', false)
  if (d) {
    return new Date(d).toString()
  }
  return ''
})
</script>
