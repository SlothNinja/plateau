  <template>
    <v-card color='green'>
      <v-container class='fill-height'>

        <v-row>
          <v-col cols='12'>
            <Bids :game='game' />
          </v-col>
        </v-row>

        <v-row v-if='showForm'>
          <v-col cols='12'>
              <BidForm />
          </v-col>
        </v-row>

        <v-row>
          <v-col cols='12'>
            <Trick v-if="showTrick" />
          </v-col>
        </v-row>
      </v-container>
    </v-card>
  </template>

<script setup>
// components
import Bids from '@/components/Game/Bids.vue'
import BidForm from '@/components/Game/BidForm.vue'

// vue
import { computed, inject } from 'vue'

// lodash
import _get from 'lodash/get'

// composables
import { cuKey, gameKey } from '@/composables/keys.js'
import { useIsCP } from '@/composables/player.js'

const game = inject(gameKey)
const cu = inject(cuKey)

const isCP = computed(() => (useIsCP(game, cu)))

const showForm = computed(() => {
  const phase = _get(game, 'value.header.phase', '')
  return isCP.value && (phase == 'bid' || phase == 'increase objective')
})

const showTrick = computed(() => {
  const phase = _get(game, 'value.header.phase', '')
  return phase == 'card play'
})

</script>
