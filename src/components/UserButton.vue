<template>
  <div>
    <v-avatar color='black' :size='size' class='ma-1'>
      <v-avatar :image='path.value' :size='0.9*size'></v-avatar>
    </v-avatar>
    <slot>
    User Name
    </slot>
  </div>
</template>

<script setup>

  import  { useGravatar } from '@/composables/gravatar.js'
  import { computed } from 'vue'

  const props = defineProps({
    color: { type: String, default: 'black' },
    user: { type: Object, required: true },
    size: { type: Number, required: true },
  })

  const path = computed( () => useGravatar(props.user.emailHash ))
</script>
