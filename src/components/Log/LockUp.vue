<template>
  <li class='my-1'>
    Locked Up Ward {{message.ward}}.
  </li>
</template>

<script>

export default {
  name: 'sn-log-lock-up-msg',
  props: [ 'message', 'game' ],
}
</script>
