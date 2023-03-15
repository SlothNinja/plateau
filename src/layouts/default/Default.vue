<script setup>
  import DefaultToolBar from './ToolBar.vue'
  import DefaultNavDrawer from './NavDrawer.vue'
  import DefaultView from './View.vue'
  import DefaultFooter from './Footer.vue'
  import { ref } from 'vue'

  const nav = ref(false)

</script>

<template>
  <v-app>
    <default-tool-bar v-model='nav' />

    <default-view />

    <default-footer />

    <default-nav-drawer v-model='nav' />
  </v-app>
</template>
