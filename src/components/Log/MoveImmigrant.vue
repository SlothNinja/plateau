<template>
  <li class='my-1'>
    Used Precinct Chairman Office to move
    <sn-cube
      :kind='message.immigrant'
      class='d-inline-flex'
      >
      1
    </sn-cube>
    from ward {{message.fromWard}} to ward {{message.toWard}}.
  </li>
</template>

<script>
  import Cube from '@/components/game/Cube'

  export default {
    name: 'sn-log-move-immigrant-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-cube': Cube,
    }
  }
</script>
