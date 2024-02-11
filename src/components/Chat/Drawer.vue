<template>
  <v-navigation-drawer
      v-model="open"
      location='right' 
      width='500'
      >
      <div class='overflow-auto fill-height flex-column d-flex justify-space-between'>
        <v-toolbar color='green' class='text-subtitle-1'>
          <v-toolbar-title>Chat</v-toolbar-title>
        </v-toolbar>
        <v-sheet class='align-start'>
          <Message
              class='mb-1 mx-1'
              v-for='(message, index) in sorted'
              :message='message'
              >
          </Message>
        </v-sheet>

        <v-sheet class='align-end'>
          <v-textarea
              ref='chatbox'
              auto-grow
              color='green'
              label='Message'
              placeholder="Type Message.  Press 'Enter' Key To Send."
              v-model='message.text'
              rows=1
              hide-details
              clearable
              v-on:keyup.enter='send'
              >
              <template v-slot:prepend-inner>
                <v-btn v-if='canSend' color='green' size='small' density='comfortable' icon='mdi-send' @click='send'></v-btn>
              </template>
          </v-textarea>
        </v-sheet>
      </div>
  </v-navigation-drawer>
</template>

<script setup>
import Message from '@/components/Chat/Message'
import { computed, inject, ref, watch, unref, nextTick, onMounted } from 'vue'
import { cuKey, gameKey, snackKey } from '@/snvue/composables/keys.js'
import { useDebouncedRef } from '@/composables/debouncedRef'
import { doc, collection } from 'firebase/firestore'
import { useRoute } from 'vue-router'
import { db } from '@/composables/firebase'
import { usePut } from '@/snvue/composables/fetch'
import { useFirestore } from '@vueuse/firebase/useFirestore'

import _get from 'lodash/get'
import _sortBy from 'lodash/sortBy'
import _isEmpty from 'lodash/isEmpty'
import _filter from 'lodash/filter'
import _includes from 'lodash/includes'
import _size from 'lodash/size'
import _map from 'lodash/map'

// const props = defineProps(['modelValue'])
// const emit = defineEmits(['update:modelValue', 'unread'])
const emit = defineEmits(['unread'])

const route = useRoute()
const cu = inject(cuKey)


const gid = computed(() => _get(unref(route), 'params.id', ''))
const messagesRef = computed(() => collection(db, 'Game', unref(gid), 'Messages'))
const messages = useFirestore(messagesRef)
const sorted = computed(
  () => {
    if (_isEmpty(unref(messages))) {
      return []
    }
    return _sortBy(unref(messages), m => m.CreatedAt.toDate())
  }
)

const open = defineModel()

watch(open, () => {
  if(unref(open)) {
    scrollChatBox()
  }
})

watch(sorted, () => {
  if(unref(open)) {
    scrollChatBox()
  }
})

const game = inject(gameKey)

const message = ref( { text: '' })

const chatbox = ref(null)

function scrollChatBox(){
  if (chatbox.value) {
    nextTick(() => chatbox.value.scrollIntoView(false))
  }
}

const msgURL = computed(() => {
  if (process.env.NODE_ENV == 'development') {
    const backend = import.meta.env.VITE_PLATEAU_BACKEND
    return `${backend}sn/mlog/add/${unref(gid)}`
  }
  return `/sn/mlog/add/${unref(gid)}`
})

///////////////////////////////////////////////////////
// Put data of new invitation to server
function send () {
  let m = unref(message)
  m.creator = unref(cu)
  const { data: response } = usePut(msgURL, m).json()
  unref(message).text = ''
  watch(response, () => update(response))
}


function update(response) {
  const msg = _get(unref(response), 'Message', '')
  if (!_isEmpty(msg)) {
    updateSnackbar(msg)
  }
}

const canSend = computed(() => !_isEmpty(_get(unref(message), 'text', '')))

const cuid = computed(() => _get(unref(cu), 'ID', 0))
const unreadIDS = useDebouncedRef([], 5000)
const unreadMesssages = computed(
  () => {
    const um = _filter(unref(messages), (m) => !_includes(_get(m, 'Read', []), unref(cuid)))
    unreadIDS.value = _map(um, (m) => m.id)
    return um
  }
)

const unread = computed(() => {
  const value = _size(unref(unreadMesssages))
  emit('unread', value)
  return value
})

watch([open, unreadIDS], () => {
  if ((_size(unref(unreadIDS)) > 0) && (unref(open))) {
    const { data: response } = usePut(`/sn/mlog/updateRead/${route.params.id}`, { "Read": unref(unreadIDS) }).json()
    watch(response, () => update(response))
  }
})

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

</script>
