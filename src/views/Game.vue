<template>
  <v-container>
    <v-row>
      <v-col cols='6'>
        <GameTable />
      </v-col>
      <v-col cols='6'>
        <GameBoard />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import GameBoard from '@/components/GameBoard'
import GameTable from '@/components/GameTable'

</script>
