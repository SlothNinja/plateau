<template>
  <v-container>
    <v-card height='200'>
      <v-img height='200' cover :src="board36"></v-img>
    </v-card>
    <v-card class='mt-1'>
      <v-card-text>
      <v-row>
        <v-col class='font-weight-black'>Designer</v-col>
        <v-col>Raymond Gallardo</v-col>
      </v-row>
      <v-row>
        <v-col class='font-weight-black'>Artists</v-col> 
        <v-col>Raymond Gallardo</v-col> 
      </v-row>
      <v-row> 
        <v-col class='font-weight-black'>Publisher</v-col> 
        <v-col><router-link :to="{ name: 'WebPublished'}">Web Published</router-link></v-col>
      </v-row>
      <v-row>
        <v-col class='font-weight-black'>Year Published</v-col>
        <v-col>2021</v-col>
      </v-row>
      <v-row> 
        <v-col class='font-weight-black'>On-Line Developer</v-col> 
        <v-col>Jeff Huter</v-col> 
      </v-row> 
      <v-row> 
        <v-col class='font-weight-black'>Permission Provided By</v-col> 
        <v-col>Raymond Gallardo</v-col> 
      </v-row> 
      <v-row> 
        <v-col class='font-weight-black'>Rules (pdf)</v-col> 
        <v-col><router-link :to="{ name: 'Rules'}">English (BGG)</router-link></v-col> 
      </v-row> 
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script setup>
import board36 from '@/assets/board36.png'
</script>
