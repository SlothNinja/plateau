<template>
  <tr>
    <td class='text-left'>
      <UserButton :user="user" :size='32' >
      {{user.Name}}
      </UserButton>
    </td>
    <td class='text-center'>
      {{elo}}
    </td>
    <td class='text-center'>
      {{played}}
    </td>
    <td class='text-center'>
      {{won}}
    </td>
    <td class='text-center'>
      {{wp}}%
    </td>
  </tr>
</template>

<script setup>
// components
import UserButton from '@/components/Common/UserButton.vue'

// vue
import { computed } from 'vue'

// lodash
import _get from 'lodash/get'
import _find from 'lodash/find'

const props = defineProps({
  details: { type: Object },
  user: { type: Object },
})

const elo = computed(() => detailsFor('elo'))
const played = computed(() => detailsFor('played'))
const won = computed(() => detailsFor('won'))
const wp = computed(() => detailsFor('wp'))

function detailsFor (stat) {
  return _get(_get(props, 'details', []), stat, 0)
}
</script>
