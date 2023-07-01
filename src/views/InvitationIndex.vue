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
              <v-btn @click.stop='abort(item.raw.id)' size='x-small' rounded color='green'>Abort</v-btn>
            </template>

            <template v-slot:item.title='{ item }'>
              <v-icon class='mb-2' size='small' v-if='item.raw.Private'>mdi-lock</v-icon>{{item.raw.Title}}
            </template>

            <template v-slot:item.creator='{ item }'>
              <UserButton :user='useCreator(item.raw)' :size='size' />
            </template>

            <template v-slot:item.numPlayers='{ item }'>
              {{item.raw.NumPlayers}}
            </template>

            <template v-slot:item.hands='{ item }'>
              {{handsPerPlayer(item)}}
            </template>

            <template v-slot:item.players="{ item }">
              <UserButton class='mb-1' :user="user" :size='size' v-for='user in useUsers(item.raw)' :key='user.id' />
            </template>

            <template v-slot:item.lastUpdated="{ item }">
              {{fromNow(item.raw.UpdatedAt.toDate())}}
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
import { VDataTable } from 'vuetify/labs/VDataTable'

// Composables
import { useFetch, usePut } from '@/composables/fetch'
import { useCreator, useUsers } from '@/composables/user'
import { fromNow } from '@/composables/fromNow'
import { db } from '@/composables/firebase'

// Vue
import { computed, inject, ref, unref, watch } from 'vue'
import { useCollection, useFirestore } from 'vuefire'
import { collection, query, where } from 'firebase/firestore'

// Lodash
import _get from 'lodash/get'
import _size from 'lodash/size'
import _reverse from 'lodash/reverse'
import _sortBy from 'lodash/sortBy'
import _isEmpty from 'lodash/isEmpty'

// inject current user
import { cuKey, snackKey } from '@/composables/keys'
const cu = inject(cuKey)

const invitations = useCollection(query(collection(db, 'Invitation'), where('Status', '==', 'recruiting')))

const sorted = computed(() => _reverse(_sortBy(unref(invitations), ['UpdatedAt'])))

const expanded = ref([])

const headers = computed(
  () => {
    if (unref(cu).Admin) {
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
  const opt = JSON.parse(_get(item, 'raw.OptString', {}))
  return _get(opt, 'HandsPerPlayer', 0)
}

// Inject snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

function abort (id) {
  const { response, error } = usePut(`/sn/invitation/abort/${id}`)
  watch(response, () => update(response))
}

function update(response) {
  const msg = _get(unref(response), 'Message', '')
  if (!_isEmpty(msg)) {
    updateSnackbar(msg)
  }
}

</script>
