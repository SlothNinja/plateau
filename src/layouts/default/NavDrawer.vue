<template>
  <v-navigation-drawer v-model='value'>
    <v-list v-if='cu'>
      <v-list-item>
        <UserButton v-if='cu' :user='cu' :size='40'> {{name}} </UserButton>
      </v-list-item>
    </v-list>
    <v-divider></v-divider>
    <v-list nav>
      <v-list-item v-if='!cu' title='Login' :to="{ name: 'Login' }" prependIcon='mdi-login' ></v-list-item>
      <v-list-item v-if='cu' title='Logout' :to="{ name: 'Logout' }" prependIcon='mdi-logout' ></v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
  import { computed } from 'vue'
  import  { useFetch } from '@/composables/fetch.js'
  import UserButton from '@/components/UserButton.vue'
  import { get } from 'lodash'

  const props = defineProps(['modelValue'])
  const emit = defineEmits(['update:modelValue'])
  const { data, error } = useFetch('https://plateau.fake-slothninja.com:8091/sn/home')

  const value = computed({
    get() {
      return props.modelValue
    },
    set(value) {
      emit('update:modelValue', value)
    }
  })

  const cu = computed( () => {
    return get(data, 'value.cu', {})
  })

  const name = computed( () => {
    return get(cu, 'value.name', '')
  })
</script>
