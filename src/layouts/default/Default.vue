<template>
  <v-app>
    <DefaultToolBar @toggleNav='toggleNav' />

    <DefaultView />

    <DefaultFooter />

    <DefaultNavDrawer v-model='nav' />

    <DefaultSnack v-model:open='snackbar.open' v-model:message='snackbar.message' />
  </v-app>
</template>

<script setup>
import DefaultToolBar from '@/layouts/default/ToolBar.vue'
import DefaultNavDrawer from '@/layouts/default/NavDrawer.vue'
import DefaultView from '@/layouts/default/View.vue'
import DefaultFooter from '@/layouts/default/Footer.vue'
import DefaultSnack from '@/layouts/default/SnackBar.vue'
import { ref, provide } from 'vue'
import { snackKey } from '@/composables/keys.js'

const nav = ref(false)
const snackbar = ref({
  message: '',
  open: false,
})

provide( snackKey, snackbar )

function toggleNav() {
  nav.value = !nav.value
}


// function clearMessage () {
//   console.log('clearMessage')
//   if (_get(snackbar, 'value.open', false) == false) {
//     snackbar.value.message = ''
//   }
// }

// import { ref, inject } from 'vue'

// const snackbar = inject(snackKey)

// const open = ref(false)

</script>

