<template>
  <Avatar
      :color='color'
      :user='user'
      :size='size'
      @click="toUser"
      style="cursor:pointer"
      >
      <slot></slot>
  </Avatar>
</template>

<script setup>
import Avatar from '@/components/Avatar'
import { useRouter } from 'vue-router'
import _get from 'lodash/get'

const props = defineProps({
  color: { type: String, default: 'black' },
  user: { type: Object, required: true },
  size: { type: Number, required: true },
})

const router = useRouter()

function toUser() {
  let id = _get(props.user, 'id', -1)
  if (id != -1) {
    router.push({ name: 'User', params: { id: id } })
  }
}
</script>
