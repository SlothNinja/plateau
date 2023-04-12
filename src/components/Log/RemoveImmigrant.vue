<template>
  <li class='my-1'>
    Used Chief of Police Office to remove
    <sn-cube
      :kind='message.immigrant'
      class='d-inline-flex'
      >
      1
    </sn-cube>
    ward {{message.ward}}.
  </li>
</template>

<script>

  import Cube from '@/components/game/Cube'

  export default {
    name: 'sn-log-remove-immigrant-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-cube': Cube,
    }
  }
</script>
