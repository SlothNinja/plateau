<template>
  <div>
    Good luck.  Have fun.
    <div class="mt-2 d-flex flex-wrap justify-space-around">
        <sn-user-btn
          v-for='(player, index) in players'
          :key='index'
          :user='userFor(player.id)'
          :color='colorByPID(player.id)'
          x-small
          bottom
        >
          {{nameFor(player.id)}}
        </sn-user-btn>
    </div>
  </div>
</template>

<script>
  import Button from '@/components/lib/user/Button'
  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'

  export default {
    mixins: [ Player, Color ],
    name: 'sn-log-start-game-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button
    }
  }
</script>
