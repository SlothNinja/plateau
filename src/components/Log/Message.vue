<template>
  <component
      :is='component'
      :data='data'
      >
  </component>
</template>

<script setup>
import StartGame from '@/components/Log/StartGame'
import AnnounceDeal from '@/components/Log/AnnounceDeal'
import EndHandResults from '@/components/Log/EndHandResults'
import StartBid from '@/components/Log/StartBid'
import PlacedBid from '@/components/Log/PlacedBid'
import PassedBid from '@/components/Log/PassedBid'
import ExchangedCards from '@/components/Log/ExchangedCards'
import PassedIncObjective from '@/components/Log/PassedIncObjective'
import IncObjective from '@/components/Log/IncObjective'
import PlayedCard from '@/components/Log/PlayedCard'
import { computed, watch, ref, unref, onMounted } from 'vue'

const props = defineProps([ 'template', 'data' ])

const component = computed(() => {
  switch (props.template) {
    case "start-game":
      return StartGame
    case "announce-deal":
      return AnnounceDeal
    case "end-hand-results":
      return EndHandResults
    case "start-bid":
      return StartBid
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
    default:
      return ''
  }
})

// onMounted(() => title.value='This is title')

</script>
