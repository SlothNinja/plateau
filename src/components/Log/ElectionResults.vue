<template>
  <fragment>
    <sn-contested-results-msg
      v-if="message.contested"
      :message="message"
      :game="game"
      >
    </sn-contested-results-msg>
    <sn-uncontested-results-msg
      v-else
      :message="message"
      :game="game"
      >
    </sn-uncontested-results-msg>
  </fragment>
</template>

<script>
  import ContestedResults from "@/components/log/ContestedResults"
  import UncontestedResults from "@/components/log/UncontestedResults"

  export default {
    name: "sn-log-election-results-msg",
    props: [ "message", "game" ],
    components: {
      "sn-contested-results-msg": ContestedResults,
      "sn-uncontested-results-msg": UncontestedResults
    }
  }
</script>
