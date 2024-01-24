<template>
  <div :class='klass'>
    <Avatar
        :hash='user.EmailHash'
        :size='size'
        :variant='user.GravType'
        :color='color'
        @click="toUser"
        style="cursor:pointer"
        />
    <div class='ml-1'>
      <slot v-if='showSlot'></slot>
      <div v-else>{{user.Name}}</div>
    </div>
  </div>
</template>

<script setup>
import Avatar from '@/components/Common/Avatar'
import { useRouter } from 'vue-router'
import { computed, useSlots } from 'vue'
import _get from 'lodash/get'
import _isEmpty from 'lodash/isEmpty'
import _isFunction from 'lodash/isFunction'

const props = defineProps({
  user: { type: Object, required: true },
  size: { type: Number, required: true },
  bottom: { type: Boolean, default: false },
  color: { type: String, default: 'black' },
})

const router = useRouter()

function toUser() {
  const id = _get(props.user, 'ID', -1)
  if (props.id != -1) {
    router.push({ name: 'User', params: { id: id } })
  }
}

const klass = computed( () => {
  if (props.bottom) {
      return 'd-inline-flex align-center flex-column'
  }
  return 'd-flex align-center'
})

const slots = useSlots()

const showSlot = computed(() =>
  {
    if (_isFunction(slots.default)) {
      return (!_isEmpty(_get(slots.default(), '[0].children[0]', {})))
    }
    return false
  }
)

</script>
