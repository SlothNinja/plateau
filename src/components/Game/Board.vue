<template>
  <v-img :src="numPlayers == 2 ? board2 : board36">
    <div
        v-for='(card, index) in declarersCards'
        :key='index'
        class='dot dot-red'
        :class='cardClass(card.Rank, card.Suit)'
        >
    </div>
    <div
        v-for='(card, index) in opposersCards'
        :key='index'
        class='dot dot-blue'
        :class='cardClass(card.Rank, card.Suit)'
        >
    </div>
    <div
        v-for='(won, index) in wonTricks'
        :key='index'
        class='dot dot-blue'
        :class='trickClass(won, index)'
        >
    </div>
  </v-img>
</template>

<script setup>
import board2 from '@/assets/board2.png'
import board36 from '@/assets/board36.png'
import { computed, unref } from 'vue'
import _includes from 'lodash/includes'
import _map from 'lodash/map'
import _filter from 'lodash/filter'
import _flatMap from 'lodash/flatMap'
import _take from 'lodash/take'

const props = defineProps(['tricks', 'dTeam', 'numPlayers'])

const declarersCards = computed(() => (_flatMap(unref(filtered), (trick) => (isDeclarers(trick) ? trick.Cards : [] ))))

const opposersCards = computed(() => (_flatMap(unref(filtered), (trick) => (isDeclarers(trick) ? [] : trick.Cards ))))

const filtered = computed(() => _filter(props.tricks, (trick) => (trick.WonBy != 0)))

const wonTricks = computed(() => {
  let won = _map(unref(filtered), (trick) => (isDeclarers(trick)))
  if (unref(props.numPlayers) == 2) {
    return _take(won, 16)
  }
  return _take(won, 13)
})

function cardClass(rank, suit) {
  let klass = `${unref(rank)}-${unref(suit)}`
  if (unref(props.numPlayers) == 2) {
    return klass + '-2p'
  }
  return klass
}

function trickClass(won, index) {
  if (unref(props.numPlayers) == 2) {
    return `trick-${unref(index) + 1}-2p ${unref(won) ? 'dot-red' : 'dot-blue'}`
  }
  return `trick-${unref(index) + 1} ${unref(won) ? 'dot-red' : 'dot-blue'}`
}

function isDeclarers(trick) {
  return _includes(props.dTeam, trick.WonBy)
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
  background-color: rgb(26 35 126)

.dot-red
  background-color: rgb(213 0 0)

.one-diamonds, .one-clubs, .one-spades, .one-hearts 
  display: none

.one-diamonds-2p, .one-clubs-2p, .one-spades-2p, .one-hearts-2p 
  display: none

.two-diamonds, .two-clubs, .two-spades, .two-hearts 
  display: none

.two-diamonds-2p, .two-clubs-2p, .two-spades-2p, .two-hearts-2p 
  display: none

.three-diamonds, .three-clubs, .three-spades, .three-hearts 
  display: none

.three-diamonds-2p, .three-clubs-2p, .three-spades-2p, .three-hearts-2p 
  display: none

.four-diamonds, .four-clubs, .four-spades, .four-hearts 
  display: none

.four-diamonds-2p, .four-clubs-2p, .four-spades-2p, .four-hearts-2p 
  display: none

.five-diamonds, .five-clubs, .five-spades, .five-hearts 
  display: none

.five-diamonds-2p, .five-clubs-2p, .five-spades-2p, .five-hearts-2p 
  display: none

.six-diamonds, .six-clubs, .six-spades, .six-hearts 
  display: none

.six-diamonds-2p, .six-clubs-2p, .six-spades-2p, .six-hearts-2p 
  display: none

.seven-diamonds, .seven-clubs, .seven-spades, .seven-hearts 
  display: none

.seven-diamonds-2p, .seven-clubs-2p, .seven-spades-2p, .seven-hearts-2p 
  display: none

.eight-diamonds, .eight-clubs, .eight-spades, .eight-hearts 
  display: none

.eight-diamonds-2p, .eight-clubs-2p, .eight-spades-2p, .eight-hearts-2p 
  display: none

.nine-diamonds, .nine-clubs, .nine-spades, .nine-hearts 
  display: none

.nine-diamonds-2p, .nine-clubs-2p, .nine-spades-2p, .nine-hearts-2p 
  display: none

.ten-diamonds, .ten-clubs, .ten-spades, .ten-hearts 
  display: none

.ten-diamonds-2p, .ten-clubs-2p, .ten-spades-2p, .ten-hearts-2p 
  display: none

.four-trumps-2p, .five-trumps-2p, .six-trumps-2p
  display: none

.seven-trumps, .eight-trumps, .nine-trumps, .ten-trumps
  display: none

.seven-trumps-2p, .eight-trumps-2p, .nine-trumps-2p, .ten-trumps-2p
  display: none

.eleven-trumps, .twelve-trumps, .thirteen-trumps, .fourteen-trumps
  display: none

.eleven-trumps-2p, .twelve-trumps-2p, .thirteen-trumps-2p, .fourteen-trumps-2p
  display: none

.fifteen-trumps, .sixteen-trumps, .seventeen-trumps, .eighteen-trumps
  display: none

.fifteen-trumps-2p, .sixteen-trumps-2p, .seventeen-trumps-2p, .eighteen-trumps-2p
  display: none

.nineteen-trumps, .twenty-trumps, .nineteen-trumps-2p, .twenty-trumps-2p
  display: none

.trick-1, .trick-1-2p
  position: absolute
  top: 11.5%
  left: 58%

.trick-2, .trick-2-2p
  position: absolute
  top: 18.5%
  left: 70%

.trick-3, .trick-4-2p
  position: absolute
  top: 39%
  left: 82%

.trick-4, .trick-5-2p 
  position: absolute
  top: 53%
  left: 82%

.trick-5, .trick-6-2p
  position: absolute
  top: 74%
  left: 70%

.trick-6, .trick-7-2p
  position: absolute
  top: 81%
  left: 58%

.trick-7, .trick-9-2p
  position: absolute
  top: 81%
  left: 34%

.trick-8, .trick-10-2p
  position: absolute
  top: 74%
  left: 22%

.trick-9, .trick-11-2p
  position: absolute
  top: 53%
  left: 10%

.trick-10, .trick-12-2p
  position: absolute
  top: 39%
  left: 10%

.trick-11, .trick-14-2p
  position: absolute
  top: 18.5%
  left: 22%

.trick-12, .trick-15-2p
  position: absolute
  top: 11.5%
  left: 34%

.trick-13, .trick-16-2p
  position: absolute
  top: 46%
  left: 46%

.trick-14
  display: none

.dame-hearts, .dame-hearts-2p
  position: absolute
  top: 32%
  left: 22%

.dame-clubs, .dame-clubs-2p
  position: absolute
  top: 32%
  left: 70%

.dame-diamonds, .dame-diamonds-2p
  position: absolute
  top: 60%
  left: 70%

.dame-spades, .dame-spades-2p
  position: absolute
  top: 60%
  left: 22%

.cavalier-hearts, .cavalier-hearts-2p
  position: absolute
  top: 18.5%
  left: 46%

.cavalier-diamonds, .cavalier-diamonds-2p
  position: absolute
  top: 74%
  left: 46%

.cavalier-clubs, .cavalier-clubs-2p
  position: absolute
  top: 46%
  left: 22%

.cavalier-spades, .cavalier-spades-2p
  position: absolute
  top: 46%
  left: 70%

.valet-diamonds, .valet-diamonds-2p
  position: absolute
  top: 25%
  left: 34%

.valet-spades, .valet-spades-2p
  position: absolute
  top: 25%
  left: 58%

.valet-hearts, .valet-hearts-2p
  position: absolute
  top: 67%
  left: 58%

.valet-clubs, .valet-clubs-2p
  position: absolute
  top: 67%
  left: 34%

.roi-spades, .roi-spades-2p
  position: absolute
  top: 39%
  left: 34%

.roi-diamonds, .roi-diamonds-2p
  position: absolute
  top: 39%
  left: 58%

.roi-clubs, .roi-clubs-2p
  position: absolute
  top: 53%
  left: 58%

.roi-hearts, .roi-hearts-2p
  position: absolute
  top: 53%
  left: 34%

.one-trumps, .one-trumps-2p
  position: absolute
  top: 4.5%
  left: 46%

.two-trumps, .trick-3-2p
  position: absolute
  top: 25%
  left: 82%

.three-trumps, .two-trumps-2p
  position: absolute
  top: 67%
  left: 82%

.four-trumps, .trick-8-2p
  position: absolute
  top: 88%
  left: 46%

.five-trumps, .three-trumps-2p
  position: absolute
  top: 67%
  left: 10%

.six-trumps, .trick-13-2p
  position: absolute
  top: 25%
  left: 10%

.twentyone-trumps, .twentyone-trumps-2p
  position: absolute
  top: 60%
  left: 46%

.excuse-trumps, .excuse-trumps-2p
  position: absolute
  top: 32%
  left: 46%

</style>
