<template>
  <div>
    <h2>Final Scores</h2>
    <div
      class='my-1'
      v-for='(player, index) in players'
      :key='index'
      >
      <sn-user-btn
        :user='userFor(player.id)'
        :color='colorByPID(player.id)'
        x-small
      >
        {{nameFor(player.id)}}
      </sn-user-btn>
      scored {{player.score}}.
    </div>
    Congratulations:
    <sn-user-btn
      v-for='(winner, index) in winners'
      :key='index'
      :user='userFor(winner.id)'
      :color='colorByPID(winner.id)'
      x-small
    >
      {{nameFor(winner.id)}}
    </sn-user-btn>
  </div>
</template>

<script>
  import Button from '@/components/lib/user/Button'
  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'

  const _ = require('lodash')

  export default {
    mixins: [ Player, Color ],
    name: 'sn-announce-winners-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button
    },
    computed: {
      winners () {
        return _.map(this.message.winners, uid => this.playerByUID(uid))
      }
    }
  }
</script>
