<template>
  <v-snackbar
      bottom
      vertical
      v-model='openValue'
      >
      <div v-html="message"></div>
      <template v-slot:actions>
        <v-btn
            variant="text"
            @click="click"
            >
            Close
        </v-btn>
      </template>
  </v-snackbar>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps(['open', 'message'])
const emit = defineEmits(['update:open', 'update:message'])

const openValue = computed({
  get() {
    return props.open
  },
  set(value) {
    emit('update:open', value)
    if (value == false) {
      emit('update:message', '')
    }
  }
})

function click() {
  emit('update:open', false)
  emit('update:message', '')
}
</script>
