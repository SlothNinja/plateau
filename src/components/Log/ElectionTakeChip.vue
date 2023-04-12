<template>
  <li class='my-1'>
    <sn-user-btn
      :user='userFor(message.pid)'
      x-small
      :color='colorByPID(message.pid)'
      >
      {{nameFor(message.pid)}}
    </sn-user-btn>
    took
    <sn-chip
      :kind='message.immigrant'
      :size='24'
      class='d-inline-flex'
      >
      1
    </sn-chip>
  </li>
</template>

<script>

  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'
  import Button from '@/components/lib/user/Button'
  import Chip from '@/components/game/Chip'

  const _ = require('lodash')

  export default {
    mixins: [ Player, Color ],
    name: 'sn-log-take-chip-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button,
      'sn-chip': Chip,
    },
    computed: {
      immigrant: function () {
        return _.capitalize(this.message.immigrant)
      }
    }
  }

</script>
