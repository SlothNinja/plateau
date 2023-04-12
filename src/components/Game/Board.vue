<template>
  <v-img :src='board36'>
    <div
        v-for='(card, index) in declarersCards'
        :key='index'
        class='dot dot-red'
        :class="`${card.rank}-${card.suit}`"
        >
    </div>
    <div
        v-for='(card, index) in opposersCards'
        :key='index'
        class='dot dot-blue'
        :class="`${card.rank}-${card.suit}`"
        >
    </div>
    <div
        v-for='(won, index) in wonTricks'
        :key='index'
        class='dot dot-blue'
        :class="`trick-${index + 1} ${won ? 'dot-red' : 'dot-blue'}`"
        >
    </div>
  </v-img>
</template>

<script setup>
import board36 from '@/assets/board36.png'
import { computed } from 'vue'
import _includes from 'lodash/includes'
import _map from 'lodash/map'
import _filter from 'lodash/filter'
import _flatMap from 'lodash/flatMap'

const props = defineProps(['tricks', 'declarersTeam'])

const declarersCards = computed(() => (_flatMap(filtered.value, (trick) => (isDeclarers(trick) ? trick.cards : [] ))))

const opposersCards = computed(() => (_flatMap(filtered.value, (trick) => (isDeclarers(trick) ? [] : trick.cards ))))

const filtered = computed(() => _filter(props.tricks, (trick) => (trick.wonBy != 0)))

const wonTricks = computed(() => _map(filtered.value, (trick) => (isDeclarers(trick))))

function isDeclarers(trick) {
  return _includes(props.declarersTeam, trick.wonBy)
}

</script>

<style lang='sass'>
.dot
  height: 8%
  width: 8%
  opacity: 0.8
  border-color: black
  border-style: solid
  border-width: 5px
  border-radius: 50%
  display: inline-block

.dot-blue
  background-color: rgb(0 0 150)

.dot-red
  background-color: rgb(150 0 0)

.one-diamonds, .one-clubs, .one-spades, .one-hearts 
  display: none

.two-diamonds, .two-clubs, .two-spades, .two-hearts 
  display: none

.three-diamonds, .three-clubs, .three-spades, .three-hearts 
  display: none

.four-diamonds, .four-clubs, .four-spades, .four-hearts 
  display: none

.five-diamonds, .five-clubs, .five-spades, .five-hearts 
  display: none

.six-diamonds, .six-clubs, .six-spades, .six-hearts 
  display: none

.seven-diamonds, .seven-clubs, .seven-spades, .seven-hearts 
  display: none

.eight-diamonds, .eight-clubs, .eight-spades, .eight-hearts 
  display: none

.nine-diamonds, .nine-clubs, .nine-spades, .nine-hearts 
  display: none

.ten-diamonds, .ten-clubs, .ten-spades, .ten-hearts 
  display: none

.seven-trumps, .eight-trumps, .nine-trumps, .ten-trumps
  display: none

.eleven-trumps, .twelve-trumps, .thirteen-trumps, .fourteen-trumps
  display: none

.fifteen-trumps, .sixteen-trumps, .seventeen-trumps, .eighteen-trumps
  display: none

.nineteen-trumps, .twenty-trumps
  display: none

.trick-1
  position: absolute
  top: 11.5%
  left: 58%

.trick-2
  position: absolute
  top: 18.5%
  left: 70%

.trick-3
  position: absolute
  top: 39%
  left: 82%

.trick-4
  position: absolute
  top: 53%
  left: 82%

.trick-5
  position: absolute
  top: 74%
  left: 70%

.trick-6
  position: absolute
  top: 81%
  left: 58%

.trick-7
  position: absolute
  top: 81%
  left: 34%

.trick-8
  position: absolute
  top: 74%
  left: 22%

.trick-9
  position: absolute
  top: 53%
  left: 10%

.trick-10
  position: absolute
  top: 39%
  left: 10%

.trick-11
  position: absolute
  top: 18.5%
  left: 22%

.trick-12
  position: absolute
  top: 11.5%
  left: 34%

.trick-13
  position: absolute
  top: 46%
  left: 46%

.trick-14
  display: none

.dame-hearts
  position: absolute
  top: 32%
  left: 22%

.dame-clubs
  position: absolute
  top: 32%
  left: 70%

.dame-diamonds
  position: absolute
  top: 60%
  left: 70%

.dame-spades
  position: absolute
  top: 60%
  left: 22%

.cavalier-hearts
  position: absolute
  top: 18.5%
  left: 46%

.cavalier-diamonds
  position: absolute
  top: 74%
  left: 46%

.cavalier-clubs
  position: absolute
  top: 46%
  left: 22%

.cavalier-spades
  position: absolute
  top: 46%
  left: 70%

.valet-diamonds
  position: absolute
  top: 25%
  left: 34%

.valet-spades
  position: absolute
  top: 25%
  left: 58%

.valet-hearts
  position: absolute
  top: 67%
  left: 58%

.valet-clubs
  position: absolute
  top: 67%
  left: 34%

.roi-spades
  position: absolute
  top: 39%
  left: 34%

.roi-diamonds
  position: absolute
  top: 39%
  left: 58%

.roi-clubs
  position: absolute
  top: 53%
  left: 58%

.roi-hearts
  position: absolute
  top: 53%
  left: 34%

.one-trumps
  position: absolute
  top: 4.5%
  left: 46%

.two-trumps
  position: absolute
  top: 25%
  left: 82%

.three-trumps
  position: absolute
  top: 67%
  left: 82%

.four-trumps
  position: absolute
  top: 88%
  left: 46%

.five-trumps
  position: absolute
  top: 67%
  left: 10%

.six-trumps
  position: absolute
  top: 25%
  left: 10%

.twentyone-trumps
  position: absolute
  top: 60%
  left: 46%

.excuse-trumps
  position: absolute
  top: 32%
  left: 46%

</style>
