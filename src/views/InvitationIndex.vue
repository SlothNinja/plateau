<template>
  <v-container fluid>
    <v-card>
      <CardStamp title='Le Plateu' subtitle='Invitations' :src='board36' width='74' />
      <v-card-text>
        <v-data-table
            v-if='items'
            v-model:expanded='expanded'
            :headers="headers"
            :items="items"
            item-value='id'
            show-expand
            >
            <template v-slot:item.title='{ item }'>
              <v-icon class='mb-2' size='small' v-if='!item.raw.public'>mdi-lock</v-icon>{{item.title}}
            </template>
            <template v-slot:item.creator='{ item }'>
              <UserButton :user='useCreator(item.raw)' :size='size' />
            </template>
            <template v-slot:item.pRounds='{ item }'>
              {{item.raw.numPlayers}} : {{item.raw.handsPerPlayer}}
            </template>
            <template v-slot:item.players="{ item }">
              <UserButton class='mb-1' :user="user" :size='size' v-for='user in useUsers(item.raw)' :key='user.id' />
            </template>
            <template v-slot:expanded-row='{ columns, item }'>
              <Expansion
                  :item='item'
                  :columns='columns'
                  @update:item='updateItem'
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
import UserButton from '@/components/UserButton.vue'
import CardStamp from '@/components/CardStamp.vue'
import Expansion from '@/components/Invitation/Expansion.vue'
import { VDataTable } from 'vuetify/labs/VDataTable'

// Composables
import { useFetch } from '@/composables/fetch.js'
import { useCreator, useUsers } from '@/composables/user.js'

// Vue
import { computed, ref } from 'vue'

// Lodash
import _get from 'lodash/get'
import _findIndex from 'lodash/findIndex'
import _filter from 'lodash/filter'
import _size from 'lodash/size'

const expanded = ref([])

const { data, error } = useFetch('/sn/invitations')

const items = computed({
  get() {
    return _get(data, 'value.invitations', [])
  },
  set(newValue) {
    data.value.invitations = newValue
  }
})

const headers = ref([
  { title: '', key: 'data-table-expand' },
  { title: 'ID', key: 'id' },
  { title: 'Title', key: 'title' },
  { title: 'Creator', key: 'creator' },
  { title: 'Players:Hands', key: 'pRounds' },
  { title: 'Players', key: 'players' },
  { title: 'Last Updated', key: 'lastUpdated' },
])

const size = 32

function updateItem(item) {
  let numUsers = _size(item.userIds)
  if (numUsers == 0 || numUsers == item.numPlayers) {
    items.value = _filter(items.value, itm => itm.id != item.id)
    return
  }
  let index = _findIndex(items.value, [ 'id', item.id ])
  if (index >= 0) {
    items.value[index] = item
  }
}

</script>
