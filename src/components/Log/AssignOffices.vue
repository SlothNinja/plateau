<template>
  <div>
    <h2>
      Office Assignments
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
                  Office
                </th>
              </thead>

              <tbody>
                <tr
                    v-for="assignment in message.assignments"
                    :key="assignment.pid"
                    >

                    <td>
                          <sn-user-btn 
                               :user="userFor(assignment.pid)"
                               :color="colorByPID(assignment.pid)"
                               class="text-body-1"
                               x-small
                               bottom
                               >
                               {{nameFor(assignment.pid)}}
                          </sn-user-btn>
                    </td>
                    <td class="text-center">
                      {{assignment.office}}
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

  const _ = require('lodash')

  export default {
    mixins: [ Color, Player, Common ],
    name: 'sn-log-assign-offices-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button
    },
    methods: {
      officeFor: function (p) {
        let office = _.find(this.assignment, { pid: p.id })
        return this.officeName(office)
      }
    }
  }
</script>
