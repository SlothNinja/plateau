<template>
  <v-app>
    <DefaultToolBar @toggleNav='toggleNav'>

      <!-- History control -->
      <v-tooltip location='bottom' color='info' text='Game Log'  >
        <template v-slot:activator='{ props }' >
          <v-btn disabled v-bind='props' icon='mdi-history' @click='toggleLog' />
        </template>
      </v-tooltip>

      <!-- Chat control -->
      <v-tooltip location='bottom' color='info' text='Chat' >
        <template v-slot:activator='{ props }' >
          <v-btn disabled v-if='show' density='compact' v-bind='props' @click='chat = !chat' stacked >
            <v-badge :content='unread' >
              <v-icon>mdi-chat</v-icon>
            </v-badge>
          </v-btn>
          <v-btn v-else disabled icon='mdi-chat' v-bind='props' @click='chat = !chat' />
        </template>
      </v-tooltip>

      <Controlbar />
    </DefaultToolBar>

    <DefaultView />

    <DefaultFooter />

    <DefaultNavDrawer v-model='nav' />

    <LogDrawer v-model='log' />

    <DefaultSnack v-model:open='snackbar.open' v-model:message='snackbar.message' />
  </v-app>
</template>

<script setup>
import DefaultToolBar from '@/layouts/default/ToolBar.vue'
import DefaultNavDrawer from '@/layouts/default/NavDrawer.vue'
import DefaultView from '@/layouts/default/View.vue'
import DefaultFooter from '@/layouts/default/Footer.vue'
import DefaultSnack from '@/layouts/default/SnackBar.vue'
import Controlbar from '@/components/Game/Controlbar.vue'
import LogDrawer from '@/components/Log/Drawer.vue'
import { computed, ref, inject, provide, unref, watch } from 'vue'
import { cuKey, gameKey, snackKey, stackKey } from '@/composables/keys.js'
import { useDocument, useCollection } from 'vuefire'
import { doc, collection } from 'firebase/firestore'
import { db } from '@/composables/firebase'
import { useRoute } from 'vue-router'
// lodash
import _get from 'lodash/get'
import _size from 'lodash/size'
import _find from 'lodash/find'
import _isEmpty from 'lodash/isEmpty'

const route = useRoute()

const nav = ref(false)
const snackbar = ref({
  message: '',
  open: false,
})

const cu = inject(cuKey)

function updateSnackbar(msg) {
  snackbar.value.message = msg
  snackbar.value.open = true
}

provide( snackKey, { snackbar, updateSnackbar } )

const chat = ref(false)
const log = ref(false)
const unread = ref(0)

const show = computed(() => (unread > 0))

function toggleNav() {
  if (!nav.value) {
    log.value = false
  }
  nav.value = !nav.value
}

function toggleLog() {
  if (!log.value) {
    nav.value = false
  }
  log.value = !log.value
}

const stackSource = computed(
  () => doc(db, 'Stack', `${route.params.id}-${unref(cu).ID}`)
)
const dbStack = useDocument(stackSource)

const viewSource = computed(
  () => doc(db, 'Committed', route.params.id, 'View', `${unref(cu).ID}` )
)
const view = useDocument(viewSource)

const current = computed(() => _get(unref(dbStack), 'Current', -1000))
const cachedPath = computed(() => `${unref(current)}-${unref(cu).ID}`)
const cachedSource = computed(
  () => doc(db, 'Committed', route.params.id, 'Cached', unref(cachedPath))
)
const cached = useDocument(cachedSource)

const game = computed(() => (unref((_isEmpty(unref(cached))) ? view : cached)))
provide(gameKey, game)

const stack = computed(() => (_isEmpty(unref(dbStack))) ? _get(unref(game), 'Undo', {}) : unref(dbStack))
provide (stackKey, stack)

</script>
