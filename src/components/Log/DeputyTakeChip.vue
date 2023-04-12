<template>
  <li class='my-1'>
    Used deputy office to take
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

  import Chip from '@/components/game/Chip'

  export default {
    name: 'sn-log-deputy-take-chip-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-chip': Chip,
    },
  }
</script>
