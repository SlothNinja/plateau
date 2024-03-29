<template>
  <v-container fluid>
    <v-card>
      <CardStamp title='Le Plateu' subtitle='Invitations' :src='board36' width='74' />
      <v-card-text>
        <v-data-table
            v-if='sorted'
            v-model:expanded='expanded'
            :headers="headers"
            :items="sorted"
            item-value='id'
            show-expand
            >

            <template v-slot:item.admin='{ item }'>
              <v-btn @click.stop='abort(item.id)' size='x-small' rounded color='green'>Abort</v-btn>
            </template>

            <template v-slot:item.title='{ item }'>
              <v-icon class='mb-2' size='small' v-if='item.Private'>mdi-lock</v-icon>{{item.Title}}
            </template>

            <template v-slot:item.creator='{ item }'>
              <UserButton :user='useCreator(item)' :size='size' />
            </template>

            <template v-slot:item.numPlayers='{ item }'>
              {{item.NumPlayers}}
            </template>

            <template v-slot:item.hands='{ item }'>
              {{handsPerPlayer(item)}}
            </template>

            <template v-slot:item.players="{ item }">
              <UserButton class='mb-1' :user="user" :size='size' v-for='user in useUsers(item)' :key='user.id' />
            </template>

            <template v-slot:item.lastUpdated="{ item }">
              {{fromNow(item.UpdatedAt.toDate())}}
            </template>
            
            <template v-slot:expanded-row='{ columns, item }'>
              <Expansion
                  :item='item'
                  :columns='columns'
              />
            </template>

        </v-data-table>
        <div v-else>No Invitations</div>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
// Assets
import board36 from '@/assets/board36.png'

// Components
import UserButton from '@/components/Common/UserButton'
import CardStamp from '@/components/Common/CardStamp'
import Expansion from '@/components/Invitation/Expansion'

// Composables
import { useFetch, usePut } from '@/snvue/composables/fetch'
import { useCreator, useUsers } from '@/snvue/composables/user'
import { fromNow } from '@/composables/fromNow'
import { db } from '@/composables/firebase'

// Vue
import { computed, inject, ref, unref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { collection, query, where } from 'firebase/firestore'
import { useFirestore } from '@vueuse/firebase/useFirestore'

// Lodash
import _get from 'lodash/get'
import _size from 'lodash/size'
import _reverse from 'lodash/reverse'
import _sortBy from 'lodash/sortBy'
import _isEmpty from 'lodash/isEmpty'

// inject current user
import { cuKey, snackKey } from '@/snvue/composables/keys'
const cu = inject(cuKey)

const invitations = useFirestore(query(collection(db, 'Invitation'), where("Status", "==", "recruiting")))

const sorted = computed(() => _reverse(_sortBy(unref(invitations), ['UpdatedAt'])))

const expanded = ref([])

const headers = computed(
  () => {
    if (_get(unref(cu), 'Admin', false)) {
      return [
        { title: '', key: 'data-table-expand' },
        { title: 'ID', key: 'id' },
        { title: 'Title', key: 'title' },
        { title: 'Creator', key: 'creator' },
        { title: 'Number of Players', key: 'numPlayers' },
        { title: 'Hands Per Player', key: 'hands' },
        { title: 'Players', key: 'players' },
        { title: 'Last Updated', key: 'lastUpdated' },
        { title: 'Admin', key: 'admin' },
      ]
    }
    return [
      { title: '', key: 'data-table-expand' },
      { title: 'ID', key: 'id' },
      { title: 'Title', key: 'title' },
      { title: 'Creator', key: 'creator' },
      { title: 'Number of Players', key: 'numPlayers' },
      { title: 'Hands Per Player', key: 'hands' },
      { title: 'Players', key: 'players' },
      { title: 'Last Updated', key: 'lastUpdated' },
    ]
  }
)

const size = 32

function handsPerPlayer(item) {
  const opt = JSON.parse(_get(item, 'OptString', {}))
  return _get(opt, 'HandsPerPlayer', 0)
}

// Inject snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

const router = useRouter()
function abort (id) {
  const href = router.resolve({ name: 'InvitationAction', params: { action: 'abort', id: unref(id) }}).href
  const url = `${import.meta.env.VITE_PLATEAU_BACKEND}sn${href}`
  const { data: response, error } = usePut(url).json()
  watch(response, () => update(response))
}

function update(response) {
  const msg = _get(unref(response), 'Message', '')
  if (!_isEmpty(msg)) {
    updateSnackbar(msg)
  }
}

</script>
