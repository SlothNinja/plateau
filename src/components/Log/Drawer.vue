<template>
  <v-navigation-drawer
    v-model="drawer"
    location='right' 
    width='500'
    >
    <div class='fill-height flex-column d-flex justify-space-between' >
    <v-toolbar color='green' class='text-subtitle-1'>
      <v-toolbar-title>Game Log</v-toolbar-title>
    </v-toolbar>
    <v-card class='overflow-auto w-100 h-100' >
      <v-container
        ref='gamelog'
        id='gamelog'
        style='overflow-y: auto'
        >
        <Entry class='my-1' v-for='(entry, index) in log' :key="index" :entry='entry' />
        <div class='gamelog'></div>
      </v-container>
    </v-card>
    </div>
  </v-navigation-drawer>
</template>

<script setup>
import Entry from '@/components/Log/Entry'
import { computed, inject, ref, unref } from 'vue'
import { gameKey } from '@/composables/keys.js'
import _get from 'lodash/get'

const props = defineProps(['modelValue'])
const emit = defineEmits(['update:modelValue'])

const log = computed(() => _get(unref(game), 'Log', []))

const drawer = computed({
  get() {
    return props.modelValue
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

const game = inject(gameKey)
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
