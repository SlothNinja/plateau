<template>
  <v-app>
    <DefaultToolBar v-model='nav' />

    <DefaultView />

    <DefaultFooter />

    <DefaultNavDrawer v-model='nav' />

    <DefaultSnack />
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

const open = ref(false)

</script>

