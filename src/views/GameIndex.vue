<template>
  <v-container fluid >
    <v-card>
      <CardStamp title='Le Plateu' :subtitle="`${_capitalize(status)} Games`" :src='board36' width='74' />
      <v-card-text>
        <v-data-table
            v-if='items'
            :headers="headers"
            :items="items"
            item-value='id'
            @click:row='show'
            >
            <template v-slot:item.title='{ item }'>
              {{item.title}}
            </template>
          <template v-slot:item.creator='{ item }'>
            <UserButton :user='useCreator(item.raw)' :size='size' />
          </template>
          <template v-slot:item.pRounds='{ item }'>
            {{item.raw.numPlayers}} : {{item.raw.roundsPerPlayer}}
          </template>
          <template v-slot:item.players="{ item }">
            <UserButton class='mb-1' :user="user" :size='size' v-for='user in useUsers(item.raw)' :key='user.id' />
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
import UserButton from '@/components/UserButton.vue'
import CardStamp from '@/components/CardStamp.vue'
import { VDataTable } from 'vuetify/labs/VDataTable'

// Composables
import { useFetch } from '@/composables/fetch.js'
import { useCreator } from '@/composables/creator.js'
import { useUsers } from '@/composables/users.js'

// Vue
import { computed, ref } from 'vue'

// Lodash
import _get from 'lodash/get'
import _findIndex from 'lodash/findIndex'
import _filter from 'lodash/filter'
import _size from 'lodash/size'
import _capitalize from 'lodash/capitalize'

// Vue router
import { useRoute, useRouter } from 'vue-router'
const route = useRoute()
const router = useRouter()


// fetch game headers from server
const status = computed(() => _get(route, 'params.status', ''))
const url = computed(() => `/sn/games/${status.value}`)
const { data, error } = useFetch(url.value)

// defines items in table
const items = computed({
  get() {
    return _get(data, 'value.gheaders', [])
  },
  set(newValue) {
    data.value.invitations = newValue
  }
})

// defines table headings
const headers = ref([
  { title: 'ID', key: 'id' },
  { title: 'Title', key: 'title' },
  { title: 'Creator', key: 'creator' },
  { title: 'Players:Rounds Per', key: 'pRounds' },
  { title: 'Players', key: 'players' },
  { title: 'Last Updated', key: 'lastUpdated' },
])

// Provides size for user buttons
const size = 32

// Directs browser to selected game in table
function show(event, data) {
  let id = _get(data, 'item.raw.id', -1)
  if (id != -1) {
    router.push({ name: 'GameShow', params: { id: id } })
  }
}

</script>
