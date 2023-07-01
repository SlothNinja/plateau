<template>
  <v-card elevation='4'>
    <v-card-text>
      <v-row no-gutters>
        <v-col cols='12' class='text-center'>
          Set Your Bid Characteristics
        </v-col>
      </v-row>
      <v-row no-gutters>
        <v-col cols='9'>
        </v-col>
        <v-col cols='3' class='text-center'>
          Point Value
        </v-col>
      </v-row>
      <v-row dense v-if='showExchange'>
        <v-col cols='9' >
          <v-radio-group :disabled="phase != 'bid'" label='Card Exchange:' density='compact' v-model='exchange'>
            <v-radio :label="'Exchange (' + exValue('exchange') + ')'" value='exchange'/>
              <v-radio :label="'No Exchange (' + exValue('no exchange') + ')'" value='no exchange'/>
          </v-radio-group>
        </v-col>
        <v-col cols='3'>
          <div class='text-center'>{{exValue(exchange)}}</div>
        </v-col>
      </v-row>
      <v-row dense>
        <v-col cols='9'>
          <v-radio-group label='Objective:' density='compact' v-model='objective'>
            <v-radio :disabled="disableObjective('bridge')" :label="'Bridge (' + obValue('bridge') + ')'" value='bridge'/>
              <v-radio :disabled="disableObjective('y')" :label="'Y (' + obValue('y') + ')'" value='y'/>
                <v-radio :disabled="disableObjective('fork')" :label="'Fork (' + obValue('fork') + ')'" value='fork'/>
                  <v-radio :disabled="disableObjective('five sides')" :label="'5 sides (' + obValue('five sides') + ')'" value='five sides'/>
                    <v-radio :disabled="disableObjective('six sides')" :label="'6 sides (' + obValue('six sides') + ')'" value='six sides'/>
          </v-radio-group>
        </v-col>
        <v-col cols='3'>
          <div class='text-center'>{{obValue(objective)}}</div>
        </v-col>
      </v-row>

      <v-row dense v-if='showTeams'>
        <v-col cols='9'>
          <v-radio-group :disabled="phase != 'bid'" label='Teams:' density='compact' v-model='teams'>
            <v-radio v-if='showTrio' :label="'Trio (' + tValue('trio') + ')'" value='trio'/>
              <v-radio :label="'Duo (' + tValue('duo') + ')'" value='duo'/>
                <v-radio :label="'Solo (' + tValue('solo') + ')'" value='solo'/>
          </v-radio-group>
        </v-col>
        <v-col cols='3'>
          <div class='text-center'>{{tValue(teams)}}</div>
        </v-col>
      </v-row>

      <v-row dense>
        <v-col cols='9'>
          <v-radio-group label='Total Bid:' density='compact'></v-radio-group>
        </v-col>
        <v-col cols='3'>
          <div class='text-center'>{{bidValue(numPlayers, bid)}}</div>
        </v-col>
      </v-row>

      <v-row> 
        <v-col class='d-flex justify-space-around'>
        <v-btn v-if='canPass' color='green' size='small' @click='pass'>Pass</v-btn>
        <v-btn v-if='canAbdicate' color='green' size='small' @click='abdicate'>Abdicate</v-btn>
        <v-btn v-if='canSubmit' color='green' size='small' @click='submit'>Bid</v-btn>
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script setup>
// components
import UserButton from '@/components/Common/UserButton'

// vue
import { computed, inject, ref, unref, watch, onMounted } from 'vue'

// lodash
import _get from 'lodash/get'
import _find from 'lodash/find'
import _last from 'lodash/last'
import _first from 'lodash/first'
import _isEmpty from 'lodash/isEmpty'

// composables
import { usePut } from '@/composables/fetch.js'
import { cuKey, gameKey, snackKey } from '@/composables/keys.js'
import { useCP, useCPID, useIsCP } from '@/composables/player.js'
import { exchangeValue, objectiveValue, teamsValue, bidValue, minBid } from '@/composables/bid.js'
import { useRoute } from 'vue-router'

const route = useRoute()

// inject game and current user
const cu = inject(cuKey)
const game = inject(gameKey)

const cp = computed(() => useCP(game))
const cpid = computed(() => useCPID(game))
const isCP = computed(() => useIsCP(game, cu))

const bvalue = ref({})

const minObjective = ref('')
const minObjectiveValue = computed(() => obValue(unref(minObjective)))
const bids = computed(() => _get(unref(game), 'Bids', []))
const lastBid = computed(() => _last(unref(bids)))

const bid = computed({
  get() {
    switch (unref(phase)) {
      case 'bid':
        if (_isEmpty(unref(bvalue))) {
          if (_isEmpty(unref(lastBid))) {
            bvalue.value = minBid(unref(numPlayers))
          } else {
            bvalue.value = { ...unref(lastBid) }
          }
          bvalue.value.PID = unref(cpid)
        }
        return unref(bvalue)
      case 'increase objective':
        bvalue.value = { ...unref(lastBid) }
        bvalue.value.PID = unref(cpid)
        minObjective.value = _get(unref(lastBid), 'Objective', '')
        return unref(bvalue)
      default:
        bvalue.value = {}
        return unrefi(bvalue)
    }
  },
  set(value) {
    bvalue.value = value
  }
})

function exValue(exchange) {
  return exchangeValue({'Exchange': exchange})
}

const exchange = computed({
  get() {
    return _get(unref(bid), 'Exchange', '')
  },
  set(value) {
    bid.value.Exchange = value
  }
})

function  obValue(objective) {
  return objectiveValue({'Objective': objective})
}

const objective = computed({
  get() {
    return _get(unref(bid), 'Objective', '')
  },
  set(value) {
    bid.value.Objective = value
  }
})

const phase = computed(() => _get(unref(game), 'Phase', ''))
const numPlayers = computed(() => _get(unref(game), 'NumPlayers', 0))

const showExchange = computed(() => (unref(numPlayers) < 6))
const showTeams = computed(() => (unref(numPlayers) >= 4))
const showTrio = computed(() => (unref(numPlayers) == 6))

function tValue(teams) {
  return teamsValue(numPlayers, {'Teams': teams})
}

const teams = computed({
  get() {
    return _get(unref(bid), 'Teams', '')
  },
  set(value) {
    bid.value.Teams = value
  }
})

const canPass = computed(() => {
  return (unref(isCP) && !unref(cp).PerformedAction) &&
    ((unref(phase) == 'bid') ||
      (unref(phase) == 'increase objective'))
})

const canSubmit = computed(() => {
  return (unref(isCP) && !unref(cp).PerformedAction) &&
    (bidValue(numPlayers, bid) > bidValue(numPlayers, lastBid)) &&
    ((unref(phase) == 'bid') ||
      (unref(phase) == 'increase objective'))

})

const canAbdicate = computed(() => {
  return (unref(isCP) && !unref(cp).PerformedAction) && declarer(cp) &&
    (unref(lastBid).PID) != unref(cpid) &&
    (unref(phase) == 'increase objective')
})

function declarer(player) {
  return _first(_get(unref(game), 'DeclarersTeam', [])) == unref(player).ID
}

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

/////////////////////////////////////
// Submit bid to server
function submit() {
  let action = 'bid'
  if (unref(phase) == 'increase objective') {
    action = 'incObjective'
  }
  const { response, error } = usePut(`/sn/game/${action}/${route.params.id}`, bid)

  watch(response, () => update(response))
}

/////////////////////////////////////
// Send pass action to server
function pass() {
  const { response, error } = usePut(`/sn/game/passBid/${route.params.id}`)

  watch(response, () => update(response))
}

/////////////////////////////////////
// Send pass action to server
function abdicate() {
  const { response, error } = usePut(`/sn/game/abdicate/${route.params.id}`)

  watch(response, () => update(response))
}

function update(response) {
    const msg = _get(unref(response), 'Message', '')
    if (!_isEmpty(msg)) {
      updateSnackbar(msg)
    }
}

function disableObjective(obj) {
  return obValue(obj) < unref(minObjectiveValue)
}

</script>
