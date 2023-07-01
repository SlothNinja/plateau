<template>
  <component
    :is='component'
    :message='message'
  >
  </component>
</template>

<script setup>
import StartGame from '@/components/Log/StartGame'
import PlacedBid from '@/components/Log/PlacedBid'
import PassedBid from '@/components/Log/PassedBid'
import ExchangedCards from '@/components/Log/ExchangedCards'
import PassedIncObjective from '@/components/Log/PassedIncObjective'
import IncObjective from '@/components/Log/IncObjective'
import PlayedCard from '@/components/Log/PlayedCard'
import WonTrick from '@/components/Log/WonTrick'
import { computed } from 'vue'

const props = defineProps([ 'message' ])

const component = computed(() => {
  switch (props.message.Template) {
    case "start-game":
      return StartGame
    case "placed-bid":
      return PlacedBid
    case "passed-bid":
      return PassedBid
    case "exchanged-cards":
      return ExchangedCards
    case "passed-inc-objective":
      return PassedIncObjective
    case "inc-objective":
      return IncObjective
    case "played-card":
      return PlayedCard
    case "won-trick":
      return WonTrick
    default:
      return ''
  }
})

</script>
