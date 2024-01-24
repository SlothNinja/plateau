<template>
  <div :style='style'>
    <div v-if="text != ''" class='text-center font-weight-bold' :class="`text-${textcolor}`" style='white-space: pre'>{{text}}</div>
    <v-img :src='src' />
  </div>
</template>

<script setup>
import back from '@/assets/card_0.png'

// Clubs
import oneclubs from '@/assets/card_1.png'
import twoclubs from '@/assets/card_2.png'
import threeclubs from '@/assets/card_3.png'
import fourclubs from '@/assets/card_4.png'
import fiveclubs from '@/assets/card_5.png'
import sixclubs from '@/assets/card_6.png'
import sevenclubs from '@/assets/card_7.png'
import eightclubs from '@/assets/card_8.png'
import nineclubs from '@/assets/card_9.png'
import tenclubs from '@/assets/card_13.png'
import valetclubs from '@/assets/card_14.png'
import cavalierclubs from '@/assets/card_10.png'
import dameclubs from '@/assets/card_11.png'
import roiclubs from '@/assets/card_12.png'

// Diamonds
import onediamonds from '@/assets/card_15.png'
import twodiamonds from '@/assets/card_16.png'
import threediamonds from '@/assets/card_17.png'
import fourdiamonds from '@/assets/card_18.png'
import fivediamonds from '@/assets/card_19.png'
import sixdiamonds from '@/assets/card_20.png'
import sevendiamonds from '@/assets/card_30.png'
import eightdiamonds from '@/assets/card_40.png'
import ninediamonds from '@/assets/card_50.png'
import tendiamonds from '@/assets/card_22.png'
import valetdiamonds from '@/assets/card_23.png'
import cavalierdiamonds from '@/assets/card_60.png'
import damediamonds from '@/assets/card_70.png'
import roidiamonds from '@/assets/card_21.png'

// Hearts
import onehearts from '@/assets/card_25.png'
import twohearts from '@/assets/card_26.png'
import threehearts from '@/assets/card_27.png'
import fourhearts from '@/assets/card_28.png'
import fivehearts from '@/assets/card_29.png'
import sixhearts from '@/assets/card_31.png'
import sevenhearts from '@/assets/card_41.png'
import eighthearts from '@/assets/card_51.png'
import ninehearts from '@/assets/card_61.png'
import tenhearts from '@/assets/card_52.png'
import valethearts from '@/assets/card_62.png'
import cavalierhearts from '@/assets/card_71.png'
import damehearts from '@/assets/card_32.png'
import roihearts from '@/assets/card_42.png'

// Spades
import onespades from '@/assets/card_72.png'
import twospades from '@/assets/card_33.png'
import threespades from '@/assets/card_34.png'
import fourspades from '@/assets/card_35.png'
import fivespades from '@/assets/card_36.png'
import sixspades from '@/assets/card_37.png'
import sevenspades from '@/assets/card_38.png'
import eightspades from '@/assets/card_39.png'
import ninespades from '@/assets/card_43.png'
import tenspades from '@/assets/card_44.png'
import valetspades from '@/assets/card_54.png'
import cavalierspades from '@/assets/card_53.png'
import damespades from '@/assets/card_63.png'
import roispades from '@/assets/card_73.png'

// Trumps
import onetrumps from '@/assets/card_64.png'
import twotrumps from '@/assets/card_74.png'
import threetrumps from '@/assets/card_45.png'
import fourtrumps from '@/assets/card_46.png'
import fivetrumps from '@/assets/card_47.png'
import sixtrumps from '@/assets/card_48.png'
import seventrumps from '@/assets/card_49.png'
import eighttrumps from '@/assets/card_55.png'
import ninetrumps from '@/assets/card_65.png'
import tentrumps from '@/assets/card_75.png'
import eleventrumps from '@/assets/card_56.png'
import twelvetrumps from '@/assets/card_66.png'
import thirteentrumps from '@/assets/card_76.png'
import fourteentrumps from '@/assets/card_57.png'
import fifteentrumps from '@/assets/card_58.png'
import sixteentrumps from '@/assets/card_59.png'
import seventeentrumps from '@/assets/card_67.png'
import eighteentrumps from '@/assets/card_77.png'
import nineteentrumps from '@/assets/card_68.png'
import twentytrumps from '@/assets/card_78.png'
import twentyonetrumps from '@/assets/card_69.png'
import excusetrumps from '@/assets/card_24.png'

import { computed, ref, unref } from 'vue'
import _lowerCase from 'lodash/lowerCase'

const props = defineProps({
  suit: { type: String, require: true },
  rank: { type: String, required: true },
  width: { type: Number, required: true },
  text: { type: String, required: true },
  textcolor: { type: String, default: 'black' },
  stacked: { type: Boolean, default: false },
})

const style = computed(() => {
  let styl = `width:${props.width}px;height:${props.width*2}`
  if (props.stacked) {
    return styl + ';position:relative;top:-145px'
  }
  return styl
})

const key = computed(() => _lowerCase(`${props.rank}${props.suit}`))

const src = computed(() => {
  switch (unref(key)) {
      // clubs
    case 'oneclubs':
      return oneclubs
    case 'twoclubs':
      return twoclubs
    case 'threeclubs':
      return threeclubs
    case 'fourclubs':
      return fourclubs
    case 'fiveclubs':
      return fiveclubs
    case 'sixclubs':
      return sixclubs
    case 'sevenclubs':
      return sevenclubs
    case 'eightclubs':
      return eightclubs
    case 'nineclubs':
      return nineclubs
    case 'tenclubs':
      return tenclubs
    case 'valetclubs':
      return valetclubs
    case 'cavalierclubs':
      return cavalierclubs
    case 'dameclubs':
      return dameclubs
    case 'roiclubs':
      return roiclubs

      // diamonds
    case 'onediamonds':
      return onediamonds
    case 'twodiamonds':
      return twodiamonds
    case 'threediamonds':
      return threediamonds
    case 'fourdiamonds':
      return fourdiamonds
    case 'fivediamonds':
      return fivediamonds
    case 'sixdiamonds':
      return sixdiamonds
    case 'sevendiamonds':
      return sevendiamonds
    case 'eightdiamonds':
      return eightdiamonds
    case 'ninediamonds':
      return ninediamonds
    case 'tendiamonds':
      return tendiamonds
    case 'valetdiamonds':
      return valetdiamonds
    case 'cavalierdiamonds':
      return cavalierdiamonds
    case 'damediamonds':
      return damediamonds
    case 'roidiamonds':
      return roidiamonds

      // hearts
    case 'onehearts':
      return onehearts
    case 'twohearts':
      return twohearts
    case 'threehearts':
      return threehearts
    case 'fourhearts':
      return fourhearts
    case 'fivehearts':
      return fivehearts
    case 'sixhearts':
      return sixhearts
    case 'sevenhearts':
      return sevenhearts
    case 'eighthearts':
      return eighthearts
    case 'ninehearts':
      return ninehearts
    case 'tenhearts':
      return tenhearts
    case 'valethearts':
      return valethearts
    case 'cavalierhearts':
      return cavalierhearts
    case 'damehearts':
      return damehearts
    case 'roihearts':
      return roihearts


      // spades
    case 'onespades':
      return onespades
    case 'twospades':
      return twospades
    case 'threespades':
      return threespades
    case 'fourspades':
      return fourspades
    case 'fivespades':
      return fivespades
    case 'sixspades':
      return sixspades
    case 'sevenspades':
      return sevenspades
    case 'eightspades':
      return eightspades
    case 'ninespades':
      return ninespades
    case 'tenspades':
      return tenspades
    case 'valetspades':
      return valetspades
    case 'cavalierspades':
      return cavalierspades
    case 'damespades':
      return damespades
    case 'roispades':
      return roispades

      // trumps
    case 'onetrumps':
      return onetrumps
    case 'twotrumps':
      return twotrumps
    case 'threetrumps':
      return threetrumps
    case 'fourtrumps':
      return fourtrumps
    case 'fivetrumps':
      return fivetrumps
    case 'sixtrumps':
      return sixtrumps
    case 'seventrumps':
      return seventrumps
    case 'eighttrumps':
      return eighttrumps
    case 'ninetrumps':
      return ninetrumps
    case 'tentrumps':
      return tentrumps
    case 'eleventrumps':
      return eleventrumps
    case 'twelvetrumps':
      return twelvetrumps
    case 'thirteentrumps':
      return thirteentrumps
    case 'fourteentrumps':
      return fourteentrumps
    case 'fifteentrumps':
      return fifteentrumps
    case 'sixteentrumps':
      return sixteentrumps
    case 'seventeentrumps':
      return seventeentrumps
    case 'eighteentrumps':
      return eighteentrumps
    case 'nineteentrumps':
      return nineteentrumps
    case 'twentytrumps':
      return twentytrumps
    case 'twentyonetrumps':
      return twentyonetrumps
    case 'excusetrumps':
      return excusetrumps

    default:
      return back
  }
})

</script>

<style lang='sass'>
.card
  width:80px
</style>
