<template>
  <v-avatar :color='color' :size='size'>
    <v-avatar :image='path' :size='size-6'></v-avatar>
  </v-avatar>
</template>

<script setup>

import  { useGravatar } from '@/composables/gravatar.js'
import { computed } from 'vue'

const props = defineProps({
  color: { type: String, default: 'black' },
  hash: { type: String, required: true },
  size: { type: Number, required: true },
  variant: { type: String, default: 'monsterid' },
})

const path = computed( () => useGravatar(props.hash, props.size, props.variant ))

</script>
