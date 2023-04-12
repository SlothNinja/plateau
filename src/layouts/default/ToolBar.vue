<template>
  <v-app-bar color='green' class='d-flex' >
    <div class='d-flex w-100 justify-space-between'>
      <v-tooltip location='bottom' text='Menu' color='info' >
        <template v-slot:activator='{ props }'>
          <v-app-bar-nav-icon v-bind='props' @click="$emit('toggleNav')" ></v-app-bar-nav-icon>
        </template>
      </v-tooltip>
      <slot></slot>
      <v-sheet color='white' min-width='100' class='mx-1' >
        <v-img width='100' :src="logo100" />
      </v-sheet>
    </div>
  </v-app-bar>
</template>


<script setup>
  import { computed } from 'vue'
  import logo100 from '@/assets/slothninja_logo_100.png'

  // const props = defineProps(['modelValue'])
  // const emit = defineEmits(['update:modelValue'])
  defineEmits(['toggleNav'])

  // const value = computed({
  //   get() {
  //     return props.modelValue
  //   },
  //   set(value) {
  //     emit('update:modelValue', value)
  //   }
  // })
</script>

<!--
<script>
  import CurrentUser from '@/components/lib/mixins/CurrentUser'
  import ImagePaths from '@/components/lib/mixins/ImagePaths'

  export default {
    mixins: [ CurrentUser, ImagePaths ],
    name: 'sn-toolbar',
    props: [ 'value' ]
  }
</script>
-->

