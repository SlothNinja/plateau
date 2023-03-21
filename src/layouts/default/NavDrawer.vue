<template>
  <v-navigation-drawer v-model='value'>
    <v-list v-if='cu'>
      <v-list-item>
        <UserButton v-if='cu' :user='cu' :size='32'> {{name}} </UserButton>
      </v-list-item>
    </v-list>
    <v-divider></v-divider>
    <v-list v-model:opened="open">
      <v-list-item prepend-icon="mdi-home" title="Home" :to="{ name: 'Home'}" ></v-list-item>
      <v-list-item v-if='cu' v-bind="props" prepend-icon="mdi-pencil" title="Create" :to="{ name: 'NewInvitation' }" > </v-list-item>
      <v-list-item v-if='cu' v-bind="props" prepend-icon="mdi-plus" title="Join" :to="{ name: 'InvitationIndex' }" > </v-list-item>
      <v-list-item v-if='cu' v-bind="props" prepend-icon="mdi-play" title="Play" :to="{ name: 'GameIndex', params: { status: 'running' } }" > </v-list-item>
      <v-divider></v-divider>
      <v-list-item v-if='!cu' title='Login' :to="{ name: 'Login' }" prependIcon='mdi-login' ></v-list-item>
      <v-list-item v-if='cu' title='Logout' :to="{ name: 'Logout' }" prependIcon='mdi-logout' ></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
import { computed, ref, onMounted, inject } from 'vue'
import UserButton from '@/components/UserButton.vue'
import { cuKey } from '@/composables/keys.js'
import { get } from 'lodash'

const props = defineProps(['modelValue', 'cu'])
const emit = defineEmits(['update:modelValue'])

const value = computed({
  get() {
    return props.modelValue
  },
  set(value) {
    emit('update:modelValue', value)
  }
})

const cu = inject(cuKey)

const name = computed( () => {
  return get(cu, 'value.name', '')
})

const open = ref( ['Create'] )

</script>
