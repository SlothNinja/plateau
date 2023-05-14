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

        <v-text-field
            label="Title"
            v-model="invitation.Title"
            >
        </v-text-field>

          <v-select
              label="Number of Players"
              :items="[ 2, 3, 4, 5, 6 ]"
              v-model="invitation.NumPlayers"
              >
          </v-select> 

            <v-select
                label="Hands per Player"
                :items="[ 1, 2, 3, 4, 5 ]"
                v-model="invitation.HandsPerPlayer"
                >
            </v-select> 

              <v-text-field
                  disabled
                  label='Password'
                  v-model='invitation.Password'
                  :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
                  :type="show ? 'text' : 'password'"
                  :placeholder="passwordMessage"
                  :hint="passwordMessage"
                  persistent-hint
                  clearable
                  @click:append="show = !show"
                  >
              </v-text-field>

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
import  { useFetch, usePut } from '@/composables/fetch.js'
import  { snackKey } from '@/composables/keys.js'

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
const { data, error } = useFetch('/sn/invitation/new')

/////////////////////////////////////////////////////
// Watch for data promise to resolve from useFetch
// update invitation to received default values
// update snackbar message based on any received message
watch(data, () => update(data))

///////////////////////////////////////////////////////
// Put data of new invitation to server
function putData () {
  const { response, error } = usePut('/sn/invitation/new', invitation)
  watch(response, () => update(response))
}

function update(data) {
  invitation.value = _get(unref(data), 'Invitation', {})
  const opt = JSON.parse(_get(unref(invitation), 'OptString', {}))
  invitation.value.HandsPerPlayer = _get(opt, 'HandsPerPlayer', 0)
  const message = _get(unref(data), 'Message', '')
  if (!_isEmpty(message)) {
    updateSnackbar(message)
  }
}

//////////////////////////////////////
// Snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

</script>
