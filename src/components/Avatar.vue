<template>
  <div class='d-flex align-center'>
    <v-avatar color='black' :size='size'>
      <v-avatar :image='path' :size='size-6'></v-avatar>
    </v-avatar>
    <div class='ml-1'>
      <slot></slot>
    </div>
  </div>
</template>

<script setup>

import  { useGravatar } from '@/composables/gravatar.js'
import { computed } from 'vue'
import { get } from 'lodash'

const props = defineProps({
  color: { type: String, default: 'black' },
  user: { type: Object, required: true },
  size: { type: Number, required: true },
})

const path = computed( () => useGravatar(props.user.emailHash, props.size, props.user.gravType ))

</script>
