<template>
  <v-sheet elevation='4' rounded >
    <div class='d-flex flex-column align-center justify-center h-100'>
      <UserButton :user='user' :size='64' variant='bottom' />
      <div>
        Score: 0
      </div>
    </div>
  </v-sheet>
</template>

<script setup>
import UserButton from '@/components/UserButton.vue'

const props = defineProps([ 'user' ])
</script>
