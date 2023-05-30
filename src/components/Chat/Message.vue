<template>
  <v-card :id='id'>
    <v-system-bar
      color='green'
      class='white--text'
    >
      <sn-user-btn
        :user='creator'
        :color='colorByUser(creator)'
        x-small
      >
      {{creator.name}}
      </sn-user-btn>


      <v-spacer></v-spacer>

      <v-tooltip bottom color='info'>
        <template v-slot:activator='{ on }'>
          <v-btn icon>
            <v-icon color='white' v-on='on'>help</v-icon>
          </v-btn>
        </template>
        <span>Help</span>
      </v-tooltip>
      <div>{{message.id}}</div>
    </v-system-bar>

    <v-card-text>
      <div>
      {{message.text}}
      </div>

      <v-divider></v-divider>

      <div class='caption'>
        Created: {{createdAt}}
      </div>
    </v-card-text>

  </v-card>
</template>

<script setup>
  import UserButton from '@/components/Common/UserButton.vue'

  // export default {
  //   mixins: [ Color ],
  //   name: 'sn-message',
  //   props: [ 'message', 'id', 'game' ],
  //   components: {
  //     'sn-user-btn': UserButton
  //   },
  //   computed: {
  //     creator: function() {
  //       return {
  //         id: this.message.creatorId,
  //         name: this.message.creatorName,
  //         emailHash: this.message.creatorEmailHash,
  //         gravType: this.message.creatorGravType
  //       }
  //     },
  //     createdAt () {
  //       var d = _.get(this.message, 'createdAt', false)
  //       if (d) {
  //         return new Date(d).toString()
  //       }
  //       return ''
  //     },
  //   }
  // }
</script>

