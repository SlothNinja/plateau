<template>
  <v-container>
    <v-row>
      <v-col cols='6'>
        <v-img :src='board36' ></v-img>
      </v-col>
      <v-col cols='6'>
        <Greeting />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import Greeting from '@/components/Home/Greeting.vue'
import board36 from '@/assets/board36.png'
</script>
