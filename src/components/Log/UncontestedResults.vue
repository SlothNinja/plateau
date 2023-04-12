<template>
  <div>
    <li class='my-1' v-if='message.winnerID'>
      <sn-user-btn
        :user='userFor(message.winnerID)'
        x-small
        :color='colorByPID(message.winnerID)'
        >
        {{nameFor(message.winnerID)}}
      </sn-user-btn>
      won the uncontested election in ward {{message.wardID}}.
    </li>
    <li v-else class='my-1'>
        No one won the uncontested election in ward {{message.wardID}}.
    </li>
  </div>
</template>

<script>
  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'
  import Button from '@/components/lib/user/Button'

  export default {
    mixins: [ Player, Color ],
    name: 'sn-uncontested-results-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button,
    }
  }
</script>
