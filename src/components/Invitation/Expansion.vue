<template>
  <tr>
    <td :colspan='length'>
          <v-table>
            <thead>
              <tr>
                <th class='text-left'>
                  Player
                </th>
                <th class='text-center'>
                  ELO
                </th>
                <th class='text-center'>
                  Played
                </th>
                <th class='text-center'>
                  Won
                </th>
                <th class='text-center'>
                  Win%
                </th>
              </tr>
            </thead>
            <tbody>
              <ExpansionRow v-for='user in users' :details='detailsFor(user.id)' :user='user' :key='user.id' />
              <ExpansionRow v-if='notJoined' :details='detailsFor(cuid)' :user='cu' :key='cuid' />
            </tbody>
          </v-table>
    </td>
  </tr>
  <tr>
    <td v-if='notJoined && !privat' :colspan='length'>
      <v-btn 
           size='small'
           rounded
           @click.native="action({ action: 'accept', item: item })"
           color='green'
           dark
           >
           Accept
      </v-btn>
    </td>
    <td v-if='notJoined && privat' :colspan='length'>
      <div class='d-flex align-center ma-4'>
         <v-btn 
             size='small'
             rounded
             @click.native="action({ action: 'accept', password: password, item: item })"
             color='green'
             dark
             >
             Accept
         </v-btn>
           <div class='mx-4 w-50'>
             <v-text-field
                 density='compact'
                 variant='underlined'
                 hide-details
                 v-model='password'
                 :append-icon="show ? 'mdi-eye' : 'mdi-eye-off'"
                 :type="show ? 'text' : 'password'"
                 label='Password'
                 placeholder='Enter Password'
                 clearable
                 autofocus
                 dense
                 outlined
                 rounded
                 @click:append="show = !show"
                 @keyup.enter="action({ action: 'accept', password: password, item: item })"
                 >
             </v-text-field>
           </div>
           </div>
    </td>
    <td v-if='!notJoined' :colspan='length'>
      <v-btn 
           size='small'
           rounded
           @click.native="action({ action: 'drop', item: item })"
           color='green'
           dark
           >
           Drop
      </v-btn>
    </td>
  </tr>
</template>

<script setup>
// components
import ExpansionRow from '@/components/Invitation/ExpansionRow'

// lodash
import _get from 'lodash/get'
import _find from 'lodash/find'
import _filter from 'lodash/filter'
import _includes from 'lodash/includes'
import _isEmpty from 'lodash/isEmpty'

// Vue
import { computed, ref, unref, inject, watch } from 'vue'

// composables
import { useFetch, usePut } from '@/composables/fetch.js'
import { useCreator, useUsers } from '@/composables/user.js'
import { cuKey, snackKey } from '@/composables/keys.js'

// Props
const props = defineProps({
  columns: { type: Object },
  item: { type: Object },
})

// Emit
// const emit = defineEmits(['update:item'])

///////////////////////////////////////
// Current User
const cu = inject(cuKey)
const cuid = computed(() => _get(unref(cu), 'ID', 0))

// Password related values
const password = ref('')
const show = ref(false)

// Display variable used to specify width of column
const length = computed(() => _get(props, 'columns.length', 1))

// Fetch player details for item
const id = ref(_get(props, 'item.raw.id', false))
const { data, error } = useFetch(`/sn/invitation/details/${unref(id)}`)
const details = computed(() => _get(unref(data), 'details', []))

// Pull details for specific user from the fetched details
function detailsFor (uid) {
  return _find(unref(details), { 'id': uid })
}

// Create creator and user objects from item
const creator = computed(() => useCreator(_get(props, 'item.raw', {})))
const users = computed(() => useUsers(_get(props, 'item.raw', [])))

// Indicates whether the current user has joined the item
const notJoined = computed(() => !_includes(_get(props, 'item.raw.UserIDS', []), unref(cuid)))

// Indicates whether the item is public or password protected
// const publick = computed(() => _get(props, 'item.raw.public', false))
const privat = computed(() => _get(props, 'item.raw.Private', false))


// Inject snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

// Accept or drop from invitation
function action(obj) {
  const action = _get(obj, 'action', '')
  const id = _get(obj, 'item.raw.id', 0)
  const pword = _get(obj, 'password', '')
  const { response, error } = usePut(`/sn/invitation/${action}/${id}`, { password: pword })

  // Wait for response data from server and update snackbar and clear password
  watch(response, () => {
    const message = _get(unref(response), 'Message', '')
    if (!_isEmpty(message)) {
      updateSnackbar(message)
    }
    password.value = ''
  })
}
</script>
