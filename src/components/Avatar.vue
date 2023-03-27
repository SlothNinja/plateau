<template>
  <div>
    <div :class='klass'>
      <v-avatar color='black' :size='size'>
        <v-avatar :image='path' :size='size-6'></v-avatar>
      </v-avatar>
      <div class='ml-1'>
        <div>{{user.name}}</div>
      </div>
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
  variant: { type: String, default: 'right' },
})

const path = computed( () => useGravatar(props.user.emailHash, props.size, props.user.gravType ))
const klass = computed( () => {
  switch (props.variant) {
    case 'bottom':
      return 'd-inline-flex align-center flex-column'
    default:
      return 'd-flex align-center'
  }
})

</script>
