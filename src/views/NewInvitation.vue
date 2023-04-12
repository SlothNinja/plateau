<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <Form />
      </v-col>
      <v-col>
        <Permissions />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import Permissions from '@/components/Invitation/Permissions.vue'
import Form from '@/components/Invitation/Form.vue'

</script>
