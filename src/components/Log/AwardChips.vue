<template>
  <div>
    <h2>
      Awarded Favor Chips
    </h2>
    <v-row>
      <v-col>
        <v-simple-table
            dense
            >
            <template v-slot:default>

              <thead>
                <th class="text-center">
                  Player
                </th>
                <th class="text-center">
                  Favor Chips
                </th>
              </thead>

              <tbody>
                <tr
                    v-for="(winner, pid) in message.winners"
                    :key="pid"
                    >

                    <td>
                      <sn-user-btn 
                           :user="userFor(pid)"
                           :color="colorByPID(pid)"
                           class="text-body-1"
                           x-small
                           bottom
                           >
                           {{nameFor(pid)}}
                      </sn-user-btn>
                    </td>
                    <td>
                      <sn-immigrant-chips
                          :chips="winner"
                          :size="24"
                          class="d-flex justify-center align-center"
                          >
                      </sn-immigrant-chips>
                    </td>
                </tr>
              </tbody>
            </template>
        </v-simple-table>
      </v-col>
    </v-row>
  </div>
</template>

<script>
  import Button from '@/components/lib/user/Button'
  import Color from '@/components/mixins/Color'
  import Player from '@/components/lib/mixins/Player'
  import Common from '@/components/mixins/Common'
  import ImmigrantChips from '@/components/game/ImmigrantChips'

  export default {
    mixins: [ Color, Player, Common ],
    name: 'sn-log-award-chips-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-immigrant-chips': ImmigrantChips,
      'sn-user-btn': Button
    }
  }
</script>
