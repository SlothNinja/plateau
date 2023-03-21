<template>
  <v-container fluid>
    <v-row>
      <v-col>
        <InvitationForm />
      </v-col>
      <v-col>
        <Permissions />
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import Permissions from '@/components/Permissions.vue'
import InvitationForm from '@/components/InvitationForm.vue'

</script>
