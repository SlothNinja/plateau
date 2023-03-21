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
            v-model="invitation.title"
            >
        </v-text-field>

          <v-select
              label="Number of Players"
              :items="[ 2, 3, 4, 5, 6 ]"
              v-model="invitation.numPlayers"
              >
          </v-select> 

            <v-select
                label="Rounds per Player"
                :items="[ 1, 2, 3, 4, 5 ]"
                v-model="invitation.roundsPerPlayer"
                >
            </v-select> 

              <v-text-field
                  label='Password'
                  v-model='invitation.password'
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
import { ref, watch, inject } from 'vue'

/////////////////////////////////////////////////////
// Composables
import  { useFetch } from '@/composables/fetch.js'
import  { usePut } from '@/composables/put.js'
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
watch( data, () => {
  invitation.value = _get(data, 'value.invitation', {})
  message.value = _get(data, 'value.message', '')
})

///////////////////////////////////////////////////////
// Put data of new invitation to server
function putData () {
  const { response, error } = usePut('/sn/invitation/new', invitation)

  watch( response, () => {
    invitation.value = _get(response, 'value.invitation', {})
    message.value = _get(response, 'value.message', '')
  })
}

//////////////////////////////////////
// Snackbar
const message = ref()
const snackbar = inject(snackKey)

watch(message, (newMessage) => {
  if (!_isEmpty(newMessage)) {
    snackbar.value.message = newMessage
    snackbar.value.open = true
  }
})

</script>
