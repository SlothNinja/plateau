<template>
  <v-container class='h-100 w-100'>
    <v-card class='h-100'>
      <v-card-title>
        Le Plateau
      </v-card-title>
      <v-card-subtitle>
        New Invitation
      </v-card-subtitle>
      <v-card-text v-if='invitation'>
        <form>
        <v-text-field
            label="Title"
            v-model="invitation.Title"
            autocomplete='title'
            >
        </v-text-field>

          <v-select
              label="Number of Players"
              :items="[ 2, 3, 4, 5, 6 ]"
              v-model="invitation.NumPlayers"
              autocomplete='number-of-players'
              >
          </v-select> 

            <v-select
                label="Hands per Player"
                :items="[ 1, 2, 3, 4, 5 ]"
                v-model="invitation.HandsPerPlayer"
                autocomplete='hands-per-player'
                >
            </v-select> 

              <v-text-field
                  label='Password'
                  v-model='invitation.Password'
                  :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
                  :type="show ? 'text' : 'password'"
                  :placeholder="passwordMessage"
                  :hint="passwordMessage"
                  autocomplete='new-password'
                  persistent-hint
                  clearable
                  @click:append="show = !show"
                  >
              </v-text-field>
        </form>

                <v-btn class="mt-3" color='green' dark @click="putData">Submit</v-btn>

      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
////////////////////////////////////////////////////
// Vue
import { ref, unref, watch, inject } from 'vue'

/////////////////////////////////////////////////////
// Composables
import  { useFetch, usePut } from '@/snvue/composables/fetch.js'
import  { snackKey } from '@/snvue/composables/keys.js'

////////////////////////////////////////////////////
// lodash
import  _get from 'lodash/get'
import  _isEmpty from 'lodash/isEmpty'

////////////////////////////////////////////////////
// Toggle state of password 'show'
const show = ref( false )
const passwordMessage = "Leave empty for Public Game."

const invitation = ref()

//////////////////////////////////////////////////////
// Fetch invitation default values from server
const fetchURL = `${import.meta.env.VITE_PLATEAU_BACKEND}sn/invitation/new`
const { data: fetchResponse } = useFetch(fetchURL).json()

/////////////////////////////////////////////////////
// Watch for data promise to resolve from useFetch
// update invitation to received default values
// update snackbar message based on any received message
watch(fetchResponse, () => update(fetchResponse))

///////////////////////////////////////////////////////
// Put data of new invitation to server
function putData () {
  invitation.value.OptString = `{ "HandsPerPlayer": ${invitation.value.HandsPerPlayer} }`
  invitation.value.Type = 'plateau'
  const putURL = `${import.meta.env.VITE_PLATEAU_BACKEND}sn/invitation/new`
  const { data: response } = usePut(putURL, invitation).json()
  watch(response, () => update(response))
}

function update(response) {
  invitation.value = _get(unref(response), 'Invitation', {})
  if (!_isEmpty(unref(invitation))) {
    invitation.value.NumPlayers = 2
    invitation.value.HandsPerPlayer = 1
  }

  // const opt = JSON.parse(_get(unref(invitation), 'OptString', {}))
  // invitation.value.HandsPerPlayer = _get(opt, 'HandsPerPlayer', 0)
  const message = _get(unref(response), 'Message', '')
  if (!_isEmpty(message)) {
    updateSnackbar(message)
  }
}

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

</script>
