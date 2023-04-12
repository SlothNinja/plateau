<template>
  <v-card>
    <v-toolbar height='1em' flat color='green'>
      <v-toolbar-title class='text-subtitle-2'>
      Hand: {{entry.round}}
      </v-toolbar-title>
    </v-toolbar>
      <v-card-text>
        {{entry.messages}}
      </v-card-text>

    <!--
    <div :class='outerClass'>
      <div v-if='player' class='d-flex align-center ma-3'>
        <sn-user-btn
          :user='userFor(player.id)'
          medium
          :color='colorByPID(player.id)'
          bottom
        >
        {{nameFor(player.id)}}
        </sn-user-btn>
      </div>
      <div :class='innerClass'>
        <ul>
          <sn-log-message
            v-for='(message, index) in entry.messages'
            :key='index'
            :message='message'
            :game='game'
            :pid='pid'
          >
          </sn-log-message>
        </ul>
      </div>
    </div>
    <v-divider></v-divider>
    <div class='d-flex text-center caption'>
      <v-col>{{updatedAt}}</v-col>
    </div>
    -->
  </v-card>
</template>

<script setup>
import { computed, inject } from 'vue'
import { gameKey } from '@/composables/keys.js'
import _get from 'lodash/get'
import _first from 'lodash/first'

const props = defineProps(['entry'])

const { game, updateGame } = inject(gameKey)
const messages = computed(() => (_get(props, 'entry.messages', [])))
const first = computed(() => (_first(messages.value)))
const round = computed(() => _get(props, 'entry.round', 0))
const title = computed(() => (`Hand: ${round.value}`))
//   import Message from '@/components/log/Message'
//   import Button from '@/components/lib/user/Button'
//   import Color from '@/components/mixins/Color'
//   import Player from '@/components/lib/mixins/Player'
//   import Common from '@/components/mixins/Common'
// 
//   const _ = require('lodash')
// 
//   export default {
//     mixins: [ Color, Player, Common ],
//     name: 'sn-log-entry',
//     props: [ 'entry', 'game' ],
//     components: {
//       'sn-log-message': Message,
//       'sn-user-btn': Button
//     },
//     computed: {
//       outerClass: function () {
//         if (this.player) {
//           return 'd-flex ma-3'
//         }
//         return 'ma-3'
//       },
//       innerClass: function () {
//         if (this.player) {
//           return 'd-flex align-center ma-3'
//         }
//         return ''
//       },
//       year: function () {
//         return _.get(this.entry, 'year', 0)
//       },
//       player: function () {
//         var self = this
//         var pid = _.get(self.entry, 'pid', 0)
//         return _.find(self.players, ['id', pid])
//       },
//       pid: function () {
//         return this.player ? this.player.id : 0
//       },
//       updatedAt: function () {
//         var self = this
//         var d = _.get(self.entry, 'updatedAt', false)
//         if (d) {
//           return new Date(d).toString()
//         }
//         return ''
//       },
//       title: function () {
//         switch (this.template) {
//           case "election-results":
//             return `Term: ${this.term} Election For Ward ${this.wardID}`
//           default:
//             return `Year: ${this.year}`
//         }
//       },
//       wardID: function () {
//         return _.get(this.first, "wardID", 0)
//       },
//       first: function () {
//         return _.first(this.entry.messages)
//       },
//       template: function () {
//         return _.get(this.first, "template", "")
//       }
//     }
//   }
</script>
