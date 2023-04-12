<template>
  <li class="my-1">
    <sn-user-btn
      :user='userFor(message.pid)'
      x-small
      :color='colorByPID(message.pid)'
      >
      {{nameFor(message.pid)}}
    </sn-user-btn>
    placed
    <sn-cube
      :kind='message.immigrant'
      class='d-inline-flex'
      >
      1
    </sn-cube>
    in ward {{message.ward}}.
  </li>
</template>

<script>
  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'
  import Button from '@/components/lib/user/Button'
  import Cube from '@/components/game/Cube'

  export default {
    mixins: [ Player, Color ],
    name: 'sn-log-election-place-immigrant-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button,
      'sn-cube': Cube,
    }
  }
</script>
