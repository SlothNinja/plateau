<template>
  <li class='my-1'>
    Placed
    <sn-cube
      class='d-inline-flex ma-1'
      :kind='message.immigrant'
      >
      1
    </sn-cube>
    in ward {{message.ward}}.
  </li>
</template>

<script>
  import Cube from '@/components/game/Cube'

  export default {
    name: 'sn-log-place-immigrant-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-cube': Cube,
    }
  }
</script>
