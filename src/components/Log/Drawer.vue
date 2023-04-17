<template>
  <v-navigation-drawer
    v-model="drawer"
    location='right' 
    width='500'
    >
    <!--
    <v-card class='d-flex flex-column' height='100%' >
    -->
    <v-toolbar color='green' height='1em' class='text-subtitle-1'>
      <v-toolbar-title>Game Log</v-toolbar-title>
    </v-toolbar>
    <v-card class='w-100 h-100' >
      <v-container
        ref='gamelog'
        id='gamelog'
        style='overflow-y: auto'
        >
        <Entry class='my-1' v-for="(entry, index) in game.glog" :key="index" :entry='entry' />
        <div class='gamelog'></div>
      </v-container>
    </v-card>
  </v-navigation-drawer>
</template>

<script setup>
import Entry from '@/components/Log/Entry'
import { computed, inject, ref } from 'vue'
import { gameKey } from '@/composables/keys.js'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])

const drawer = computed({
  get() {
    return props.modelValue
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

const { game, updateGame } = inject(gameKey)
// import CurrentUser from '@/components/lib/mixins/CurrentUser'
// import Entry from '@/components/Log/Entry'
// 
// const _ = require("lodash")
// 
// export default {
//   name: 'sn-log-drawer',
//   mixins: [ CurrentUser ],
//   props: [ 'value', 'game' ],
//   components: {
//     'sn-log-entry': Entry
//   },
//   watch: {
//     drawer: function (oldValue, newValue) {
//       if (oldValue != newValue) {
//         this.scroll()
//       }
//     },
//     entries: function (oldValue, newValue) {
//       if (oldValue != newValue) {
//         this.scroll()
//       }
//     }
//   },
//   methods: {
//     scroll: function() {
//       let self = this
//       self.$nextTick(function () {
//         self.$vuetify.goTo('.gamelog', { container: self.$refs.gamelog } )
//       })
//     },
//   },
//   computed: {
//     entries: function () {
//       return _.size(this.game.glog)
//     },
//     drawer: {
//       get: function () {
//         var self = this
//         return self.value
//       },
//       set: function (value) {
//         var self = this
//         self.$emit('input', value)
//       }
//     }
//   }
// }
</script>
