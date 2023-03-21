<template>
  <div class='d-flex'>
    <div>
      <v-img class='ml-2 mt-2' :width='width' :src="src" />
    </div>
    <div>
      <v-card-title>
        {{title}}
      </v-card-title>
      <v-card-subtitle>
        {{subtitle}}
      </v-card-subtitle>
    </div>
  </div>
</template>

<script setup>
const props = defineProps([ 'title', 'subtitle', 'width', 'src' ])
</script>
