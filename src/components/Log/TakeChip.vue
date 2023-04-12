<template>
  <li class='my-1'>
    Took {{immigrant}} chip.
  </li>
</template>

<script>

const _ = require('lodash')

  export default {
    name: 'sn-log-take-chip-msg',
    props: [ 'message', 'game' ],
    computed: {
      immigrant: function () {
        return _.capitalize(this.message.immigrant)
      }
    }
  }
</script>
