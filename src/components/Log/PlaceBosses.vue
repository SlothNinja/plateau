<template>
  <li class='my-1'>
    Placed
    <sn-boss
        :color='colorByPID(pid)'
        class='d-inline-flex ma-1'
      >
      {{message.bosses}}
    </sn-boss>
    in ward {{message.ward}}.
  </li>
</template>

<script>
  import Boss from '@/components/game/Boss'
  import Color from '@/components/mixins/Color'

  export default {
    mixins: [ Color ],
    name: 'sn-log-place-bosses-msg',
    props: [ 'message', 'game', 'pid' ],
    components: {
      'sn-boss': Boss,
    },
  }
</script>
