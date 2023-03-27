<template>
  <v-app-bar color='green' class='d-flex' >
    <div class='d-flex w-100 justify-space-between'>
      <v-tooltip bottom color='info' dark>
        <template v-slot:activator='{ prop }'>
          <v-app-bar-nav-icon v-bind='prop' @click.stop="value = !value" ></v-app-bar-nav-icon>
        </template>
        <span>Menu</span>
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

  const props = defineProps(['modelValue'])
  const emit = defineEmits(['update:modelValue'])

  const value = computed({
    get() {
      return props.modelValue
    },
    set(value) {
      emit('update:modelValue', value)
    }
  })
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

