<template>
  <div>
    <h2>
      End Game Victory Points
    </h2>
        <v-simple-table
            dense
            >
            <template v-slot:default>

              <thead>
                <th class='text-center'>
                  Player
                </th>
                <th class='text-center'>
                  English
                </th>
                <th class='text-center'>
                  German
                </th>
                <th class='text-center'>
                  Irish
                </th>
                <th class='text-center'>
                  Italian
                </th>
                <th class='text-center'>
                  Slander
                </th>
              </thead>

              <tbody>
                <tr
                    v-for='(player, key) in players'
                    :key='key'
                    >

                    <td>
                      <v-row align='center'>
                        <v-col align='center'>
                          <sn-user-btn 
                               :user='userFor(player.id)'
                               :color='colorByPID(player.id)'
                               class='text-body-1'
                               x-small
                               bottom
                               >
                               {{nameFor(player.id)}}
                          </sn-user-btn>
                        </v-col>
                      </v-row>
                    </td>
                    <td class='text-center'>
                      {{vpFor(player.id, 'english')}}
                    </td>
                    <td class='text-center'>
                      {{vpFor(player.id, 'german')}}
                    </td>
                    <td class='text-center'>
                      {{vpFor(player.id, 'irish')}}
                    </td>
                    <td class='text-center'>
                      {{vpFor(player.id, 'italian')}}
                    </td>
                    <td class='text-center'>
                      {{vpFor(player.id, 'slander')}}
                    </td>
                </tr>
              </tbody>
            </template>
        </v-simple-table>
  </div>
</template>

<script>
  import Button from '@/components/lib/user/Button'
  import Color from '@/components/mixins/Color'
  import Player from '@/components/lib/mixins/Player'
  import Common from '@/components/mixins/Common'

  const _ = require('lodash')

  export default {
    mixins: [ Color, Player, Common ],
    name: 'sn-log-score-vp-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button
    },
    methods: {
      vpFor (pid, key) {
        return _.get(this.message.results[pid], key, '')
      }
    },
  }
</script>
