<template>
  <v-app>
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
          <v-btn v-if='show' density='compact' v-bind='props' @click='chat = !chat' stacked >
            <v-badge :content='unread' >
              <v-icon>mdi-chat</v-icon>
            </v-badge>
          </v-btn>
          <v-btn v-else icon='mdi-chat' v-bind='props' @click='chat = !chat' />
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
import { computed, ref, provide } from 'vue'
import { gameKey, snackKey } from '@/composables/keys.js'

const nav = ref(false)
const snackbar = ref({
  message: '',
  open: false,
})

function updateSnackbar(msg) {
  snackbar.value.message = msg
  snackbar.value.open = true
}

provide( snackKey, { snackbar, updateSnackbar } )

const game = ref({})
const chat = ref(false)
const log = ref(false)
const unread = ref(0)

const show = computed(() => (unread > 0))

function updateGame(value) {
  game.value = value
  game.selected = []
}

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

provide(gameKey, { game, updateGame })

</script>

