<template>
  <v-app>
    <ChatDrawer v-model='chat' @unread='(n) => unread = n' />

    <LogDrawer v-model='log'>
      <template v-slot:entries>
        <Message v-for='(entry, index) in entries' :key='index' :template='entry.Template' :data='entry.Data' />
      </template>
    </LogDrawer>

    <DefaultToolBar @toggleNav='toggleNav'>

      <!-- History control -->
      <v-tooltip location='bottom' color='info' text='Game Log'  >
        <template v-slot:activator='{ props }' >
          <v-btn v-bind='props' icon='mdi-history' @click='toggleLog' />
        </template>
      </v-tooltip>

      <!-- Chat control -->
      <v-tooltip location='bottom' color='info' text='Chat' >
        <template v-slot:activator='{ props }' >
          <v-btn v-if='show' density='compact' v-bind='props' @click='toggleChat' stacked >
            <v-badge :content='unread' >
              <v-icon>mdi-chat</v-icon>
            </v-badge>
          </v-btn>
          <v-btn v-else icon='mdi-chat' v-bind='props' @click='toggleChat' />
        </template>
      </v-tooltip>

      <Controlbar />

    </DefaultToolBar>

    <DefaultView />

    <DefaultFooter />

    <DefaultNavDrawer v-model='nav' />

    <DefaultSnack v-model:open='snackbar.open' v-model:message='snackbar.message' />

  </v-app>
</template>

<script setup>
import DefaultToolBar from '@/layouts/default/ToolBar'
import DefaultNavDrawer from '@/layouts/default/NavDrawer'
import DefaultView from '@/layouts/default/View'
import DefaultFooter from '@/layouts/default/Footer'
import DefaultSnack from '@/layouts/default/SnackBar'
import Controlbar from '@/components/Game/Controlbar'
import LogDrawer from '@/snvue/components/Log/Drawer'
import Message from '@/components/Log/Message'
import ChatDrawer from '@/components/Chat/Drawer'
import { computed, ref, inject, provide, unref, watch, watchEffect } from 'vue'
import { cuKey, gameKey, snackKey, stackKey } from '@/snvue/composables/keys'
import { useDocument, useCollection } from 'vuefire'
import { doc, collection } from 'firebase/firestore'
import { db } from '@/composables/firebase'
import { useRoute } from 'vue-router'
import { usePut } from '@/snvue/composables/fetch'

// lodash
import _get from 'lodash/get'
import _size from 'lodash/size'
import _find from 'lodash/find'
import _isEmpty from 'lodash/isEmpty'
import _filter from 'lodash/filter'
import _includes from 'lodash/includes'
import _map from 'lodash/map'

const route = useRoute()

const chat = ref(false)
const nav = ref(false)
const snackbar = ref({
  message: '',
  open: false,
})

const cu = inject(cuKey)
const cuid = computed(() => _get(unref(cu), 'ID', 0).toString())

function updateSnackbar(msg) {
  snackbar.value.message = msg
  snackbar.value.open = true
}

provide( snackKey, { snackbar, updateSnackbar } )

const log = ref(false)

const show = computed(() => ((unref(unread) > 0) && !unref(chat)))

function toggleNav() {
  if (!unref(nav)) {
    log.value = false
    chat.value = false
  }
  nav.value = !unref(nav)
}

const unread = ref(0)

function toggleLog() {
  if (!unref(log)) {
    nav.value = false
    chat.value = false
  }
  log.value = !unref(log)
}

function toggleChat() {
  if (!unref(chat)) {
    nav.value = false
    log.value = false
  }
  chat.value = !chat.value
}

const id = computed(() => _get(unref(route), 'params.id', 0))

const indexSource = computed(() => doc(db, 'Index', unref(id)))
const index = useDocument(indexSource)

const rev = computed(() => _get(unref(index), 'Undo.Current', -1).toString())

const stackSource = computed(
  () => doc(db, 'Stack', unref(id), 'For', unref(cuid) )
)
const dbStack = useDocument(stackSource)

const viewSource = computed(
  () => doc(db, 'Game', unref(id), 'Rev', unref(rev), 'ViewFor', unref(cuid) )
)
const view = useDocument(viewSource)

const current = computed(() => _get(unref(dbStack), 'Current', -1000).toString())
const committed = computed(() => _get(unref(dbStack), 'Committed', -1000).toString())
// const cachedPath = computed(() => `${unref(current)}-${unref(cuid)}`)
const cachedSource = computed(
  () => doc(db, 'Game', unref(id), 'Rev', unref(committed), 'CacheFor', unref(cuid), 'Rev', unref(current))
)
const cached = useDocument(cachedSource)

const game = computed(() => (unref((_isEmpty(unref(cached))) ? view : cached)))
provide(gameKey, game)

const entries = computed(() => _get(unref(game), 'Log', []))
const stack = computed(() => (_isEmpty(unref(dbStack))) ? _get(unref(game), 'Undo', {}) : unref(dbStack))
provide (stackKey, stack)

</script>
