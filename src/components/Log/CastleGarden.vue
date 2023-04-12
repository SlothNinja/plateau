<template>
  <li>
    <div v-if="message.filled">
      Placed
      <sn-immigrant-cube
        v-for="(count, nationality) in message.immigrants"
        :key='nationality'
        :kind='nationality'
        class="d-inline-flex mx-1"
        >
        {{count}}
      </sn-immigrant-cube>
      in Castle Garden.
    </div>
    <div v-else>
      Did not place cubes in the Castle Garden.
    </div>
  </li>
</template>

<script>
  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'
  import Cube from '@/components/game/Cube'

  export default {
    mixins: [ Player, Color ],
    name: 'sn-log-castle-garden-msg',
    components: {
      'sn-immigrant-cube': Cube,
    },
    props: [ 'message', 'game' ]
  }
</script>
