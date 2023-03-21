<template>
  <v-snackbar
      bottom
      vertical
      v-model='snackbar.open'
      >
      {{snackbar.message}}
       <template v-slot:actions>
        <v-btn
          variant="text"
          @click="snackbar.open = false"
        >
          Close
        </v-btn>
      </template>
  </v-snackbar>
</template>

<script setup>
import { ref, inject } from 'vue'
import { snackKey } from '@/composables/keys.js'

const snackbar = inject(snackKey)

</script>
