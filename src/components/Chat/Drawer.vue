<template>
  <v-navigation-drawer
    v-model="drawer"
    location='right' 
    width='500'
    >
    <!--
    <v-card class='d-flex flex-column' height='100%' >
    -->
    <v-toolbar color='green' height='1em' class='text-subtitle-1'>
      <v-toolbar-title>Chat</v-toolbar-title>
    </v-toolbar>
    <v-card height='100%' class='d-flex flex-column' >
    <!--
      <v-toolbar
        color='green'
        dark
        dense
        flat
        class='flex-grow-0 flex-shrink-0'
        >
        <v-toolbar-title>Chat</v-toolbar-title>

      </v-toolbar>
    -->

      <v-container
        ref='chatbox'
        fluid
        style='overflow-y: auto'
        >
        <Message
          class='my-2'
          v-for='(message, index) in messages'
          :key='index'
          :message='message'
          :id='msgId(index)'
          :game='game'
          >
        </Message>

        <!--
          <v-progress-linear
            color='green'
            indeterminate
            v-if='loading || sending'
            >
          </v-progress-linear>
        -->

            <div class='messagebox'></div>
      </v-container>
    <v-spacer></v-spacer>

    <v-divider></v-divider>

    <v-card>
      <v-card-text>

        <v-textarea
          auto-grow
          color='green'
          label='Message'
          placeholder="Type Message.  Press 'Enter' Key To Send."
          v-model='message.text'
          rows=1
          clearable
          v-on:keyup.enter='send'
          >
        </v-textarea>
        <v-btn v-if='canSend' color='green' size='small' @click='send'>Send</v-btn>

      </v-card-text>
    </v-card>

    </v-card>
<!--
    <v-card class='w-100 h-100' >
      <v-container
        ref='gamelog'
        id='gamelog'
        style='overflow-y: auto'
        >
        <Entry class='my-1' v-for='(entry, index) in log' :key="index" :entry='entry' />
        <div class='gamelog'></div>
      </v-container>
    </v-card>
-->
  </v-navigation-drawer>
</template>

<script setup>
import Message from '@/components/Chat/Message'
import { computed, inject, ref, watch, unref } from 'vue'
import { cuKey, gameKey, snackKey } from '@/composables/keys.js'
import { useDocument, useCollection } from 'vuefire'
import { doc, collection } from 'firebase/firestore'
import { useRoute } from 'vue-router'
import { db } from '@/composables/firebase'
import  { usePut } from '@/composables/fetch.js'

import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])

const route = useRoute()
const cu = inject(cuKey)

const chatSource = computed(
  () => doc(db, 'Chat', route.params.id, 'View', `${unref(cu).ID}` )
)
const chatDoc = useDocument(chatSource)

const messages = computed(() => _get(unref(chatDoc), 'Messages', []))

const drawer = computed({
  get() {
    return props.modelValue
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

const game = inject(gameKey)

const message = ref( { text: '' })

///////////////////////////////////////////////////////
// Put data of new invitation to server
function send () {
  let m = unref(message)
  m.creator = unref(cu)
  const { response, error } = usePut(`/sn/mlog/add/${route.params.id}`, m)
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

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

// import CurrentUser from '@/components/lib/mixins/CurrentUser'
// import Entry from '@/components/Log/Entry'
// 
// const _ = require("lodash")
// 
// export default {
//   name: 'sn-log-drawer',
//   mixins: [ CurrentUser ],
//   props: [ 'value', 'game' ],
//   components: {
//     'sn-log-entry': Entry
//   },
//   watch: {
//     drawer: function (oldValue, newValue) {
//       if (oldValue != newValue) {
//         this.scroll()
//       }
//     },
//     entries: function (oldValue, newValue) {
//       if (oldValue != newValue) {
//         this.scroll()
//       }
//     }
//   },
//   methods: {
//     scroll: function() {
//       let self = this
//       self.$nextTick(function () {
//         self.$vuetify.goTo('.gamelog', { container: self.$refs.gamelog } )
//       })
//     },
//   },
//   computed: {
//     entries: function () {
//       return _.size(this.game.glog)
//     },
//     drawer: {
//       get: function () {
//         var self = this
//         return self.value
//       },
//       set: function (value) {
//         var self = this
//         self.$emit('input', value)
//       }
//     }
//   }
// }
</script>
