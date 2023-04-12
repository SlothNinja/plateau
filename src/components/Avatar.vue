<template>
  <div>
    <div :class='klass'>
      <v-avatar :color='color' :size='size'>
        <v-avatar :image='path' :size='size-6'></v-avatar>
      </v-avatar>
      <div class='ml-1'>
        <slot v-if='showSlot'></slot>
        <div v-else>{{user.name}}</div>
      </div>
    </div>
  </div>
</template>

<script setup>

import  { useGravatar } from '@/composables/gravatar.js'
import { computed, useSlots } from 'vue'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'

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

const slots = useSlots()

const showSlot = computed(() => (!_isEmpty(_get(slots.default(), '[0].children[0]', {}))))

</script>
