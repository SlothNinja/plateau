<template>
  <li class='my-1'>
    Spent
    <sn-chip
      :kind='message.immigrant'
      :size='24'
      class='d-inline-flex'
      >
      {{message.count}} 
    </sn-chip>
    to slander
    <sn-user-btn
      :user='userFor(message.boss)'
      x-small
      :color='colorByPID(message.boss)'
      >
      {{nameFor(message.boss)}}
    </sn-user-btn>
    from from ward {{message.ward}}.
  </li>
</template>

<script>

  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'
  import Button from '@/components/lib/user/Button'
  import Chip from '@/components/game/Chip'

  export default {
    mixins: [ Player, Color ],
    name: 'sn-log-remove-immigrant-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button,
      'sn-chip': Chip,
    }
  }
</script>
