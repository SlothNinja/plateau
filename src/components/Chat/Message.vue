<template>
  <v-sheet elevation='4' class='px-2 py-1'>
    <div class='text-caption'>{{createdAt}}</div>
    {{message.CreatorName}}:
    {{text}}
  </v-sheet>
</template>

<script setup>

import UserButton from '@/components/Common/UserButton'
import { useCreator } from '@/composables/user'
import { computed, unref } from 'vue'
import _get from 'lodash/get'

const props = defineProps(['message'])

const creator = computed(() => useCreator(props.message))

const createdAt = computed(
  () => {
    const t = _get(props.message, 'CreatedAt', false)
    if (t) {
      return new Date(t.toDate()).toLocaleString()
    }
    return ''
  }
)

const text = computed(() => _get(props.message, 'Text', ''))
</script>
