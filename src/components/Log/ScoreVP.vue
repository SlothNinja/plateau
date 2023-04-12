<template>
  <div>
    <h2>
      Awarded Victory Points
    </h2>
        <v-simple-table
            dense
            >
            <template v-slot:default>

              <thead>
                <th class="text-center">
                  Player
                </th>
                <th class="text-center">
                  Mayor VP
                </th>
                <th class="text-center">
                  Ward VP
                </th>
                <th class="text-center">
                  Wards
                </th>
              </thead>

              <tbody>
                <tr
                    v-for="(result, pid) in message.results.playerResults"
                    :key="pid"
                    >

                    <td>
                      <v-row align="center">
                        <v-col align="center">
                          <sn-user-btn 
                               :user="userFor(pid)"
                               :color="colorByPID(pid)"
                               class="text-body-1"
                               x-small
                               bottom
                               >
                               {{nameFor(pid)}}
                          </sn-user-btn>
                        </v-col>
                      </v-row>
                    </td>
                    <td class="text-center">
                      {{ pid == message.results.mayorID ? '3' : '' }}
                    </td>
                    <td class="text-center">
                      {{result.score}}
                    </td>
                    <td class="text-center">
                      {{result.wardIDS}}
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

  export default {
    mixins: [ Color, Player, Common ],
    name: 'sn-log-score-vp-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button
    }
  }
</script>
