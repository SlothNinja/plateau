<template>
  <v-container fluid >
    <v-card>
      <CardStamp title='Le Plateu' :subtitle="`${_capitalize(status)} Games`" :src='board36' width='74' />
      <v-card-text>
        <v-data-table
            v-if='sorted'
            :headers="headers"
            :items="sorted"
            item-value='id'
            @click:row='show'
            >

            <template v-slot:item.admin='{ item }'>
              <v-btn @click.stop='abandon(item.id)' size='x-small' rounded color='green'>Abandon</v-btn>
            </template>

            <template v-slot:item.title='{ item }'>
              {{item.Title}}
            </template>

            <template v-slot:item.creator='{ item }'>
              <UserButton :user='useCreator(item)' :size='size' />
            </template>

            <template v-slot:item.pRounds='{ item }'>
              {{item.NumPlayers}} : {{handsPerPlayer(item)}}
            </template>

            <template v-slot:item.players="{ item }">
              <UserButton class='mb-1' :user="user" :size='size' v-for='user in useUsers(item)' :key='user.ID'>
              <span :class='userClass(item, user)'>{{user.Name}}</span>
              </UserButton>
            </template>

            <template v-slot:item.lastUpdated="{ item }">
              {{fromNow(item.UpdatedAt.toDate())}}
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

// Composables
import { useCreator, useUsers } from '@/composables/user'
import { fromNow } from '@/composables/fromNow'
import { usePut } from '@/composables/fetch'

// Vue
import { computed, inject, ref, unref, watch } from 'vue'
import { useCollection, useFirestore } from 'vuefire'
import { collection, query, where } from 'firebase/firestore'
import { db } from '@/composables/firebase'

// Lodash
import _get from 'lodash/get'
import _findIndex from 'lodash/findIndex'
import _filter from 'lodash/filter'
import _size from 'lodash/size'
import _capitalize from 'lodash/capitalize'
import _indexOf from 'lodash/indexOf'
import _includes from 'lodash/includes'
import _nth from 'lodash/nth'
import _isEmpty from 'lodash/isEmpty'
import _reverse from 'lodash/reverse'
import _sortBy from 'lodash/sortBy'

// inject current user
import { cuKey, snackKey } from '@/composables/keys'
const cu = inject(cuKey)
const cuid = computed(() => (_get(unref(cu), 'ID', -1)))

// Vue router
import { useRoute, useRouter } from 'vue-router'
const route = useRoute()
const router = useRouter()

const status = computed(() => _get(route, 'params.status', ''))

const items = useCollection(query(collection(db, 'Index'), where('Status', '==', unref(status))))

const sorted = computed(() => _reverse(_sortBy(unref(items), ['UpdatedAt'])))

function handsPerPlayer(item) {
  const opt = JSON.parse(_get(unref(item), 'OptString', {}))
  return _get(opt, 'HandsPerPlayer', 0)
}

// defines table headings
const headers = computed(
  () => {
    const admin = _get(unref(cu), 'Admin', false)
    if (admin) {
      return [
        { title: 'ID', key: 'id' },
        { title: 'Title', key: 'title' },
        { title: 'Creator', key: 'creator' },
        { title: 'Players:Rounds Per', key: 'pRounds' },
        { title: 'Players', key: 'players' },
        { title: 'Last Updated', key: 'lastUpdated' },
        { title: 'Admin', key: 'admin' },
      ]
    }
    return [
      { title: 'ID', key: 'id' },
      { title: 'Title', key: 'title' },
      { title: 'Creator', key: 'creator' },
      { title: 'Players:Rounds Per', key: 'pRounds' },
      { title: 'Players', key: 'players' },
      { title: 'Last Updated', key: 'lastUpdated' },
    ]
  }
)

// Provides size for user buttons
const size = 32

// Directs browser to selected game in table
function show(event, data) {
  let id = _get(data, 'item.id', -1)
  if (id != -1) {
    router.push({ name: 'GameShow', params: { id: id } })
  }
} 

function userClass(item, user) {
  const uid = _get(unref(user), 'ID', -1)

  if (item.status == 'completed') {
    return winnerClass(item, uid)
  } 
  return cpClass(item, uid)
}

function cpClass(item, uid) {
  const pid = _indexOf(_get(item, 'UserIDS', []), uid) + 1

  if (_includes(_get(item, 'CPIDS', []), pid)) {
    if (unref(cuid) == uid) {
      return 'font-weight-black text-red-darken-4'
    }
    return 'font-weight-black'
  }
  return ''
}

function winnerClass(item, uid) {
  if (_includes(_get(item, 'WinnerIDS', []), uid)) {
    if (unref(cuid) == usid) {
      return 'font-weight-black text-red-darken-4'
    }
    return 'font-weight-black'
  }
  return ''
}

// Inject snackbar
const { snackbar, updateSnackbar } = inject(snackKey)

function abandon (id) {
  const { response, error } = usePut(`/sn/game/abandon/${id}`)
  watch(response, () => update(response))
}

function update(response) {
  const msg = _get(unref(response), 'Message', '')
  if (!_isEmpty(msg)) {
    updateSnackbar(msg)
  }
}

</script>
