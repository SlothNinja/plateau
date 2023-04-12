<template>
  <div>
    <li>
      Election details:
    </li>
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
              Bosses
            </th>
            <th class="text-center">
              Favors
            </th>
            <th class="text-center">
              Total
            </th>
          </thead>

          <tbody>
              <tr v-for="p in ps" :key="p.id" >

                <td>
                  <div class="d-flex align-center justify-center my-1" >
                    <sn-user-btn 
                      :user="userFor(p.id)"
                      :color="colorByPID(p.id)"
                      class="text-body-1"
                      x-small
                      bottom
                    >
                    {{nameFor(p.id)}}
                    </sn-user-btn>
                  </div>
                </td>

                <td>
                  <div class="d-flex align-center justify-center" >
                    <sn-boss
                      :color="colorByPID(p.id)"
                      :size="24"
                      >
                      {{message.bosses[p.id]}}
                    </sn-boss>
                  </div>
                </td>

                <td>
                  <div class="d-flex align-center justify-center" >
                    <sn-immigrant-chips
                      :chips="message.playedChips[p.id]"
                      :size="24"
                      >
                    </sn-immigrant-chips>
                  </div>
                </td>

                <td>
                  <div class="d-flex align-center justify-center black--text font-weight-black">
                  {{influenceFor(p)}}
                  </div>
                </td>
              </tr>
            </tbody>
          </template>
        </v-simple-table>
      </v-col>
    </v-row>
    <v-divider class="my-1" ></v-divider>
    <li v-if="message.winnerID">
      <sn-user-btn
        :user='userFor(message.winnerID)'
        x-small
        :color='colorByPID(message.winnerID)'
        >
        {{nameFor(message.winnerID)}}
      </sn-user-btn>
      won the contested election in ward {{message.wardID}}.
    </li>
    <li v-else class="my-1">
        No one won the contested election in ward {{message.wardID}}.
    </li>
  </div>
</template>

<script>
  import Player from '@/components/lib/mixins/Player'
  import Color from '@/components/mixins/Color'
  import Button from '@/components/lib/user/Button'
  import ImmigrantChips from "@/components/game/ImmigrantChips"
  import Boss from "@/components/game/Boss"

  const _ = require("lodash")

  export default {
    mixins: [ Player, Color ],
    name: 'sn-log-election-results-msg',
    props: [ 'message', 'game' ],
    components: {
      'sn-user-btn': Button,
      "sn-boss": Boss,
      "sn-immigrant-chips": ImmigrantChips,
    },
    methods: {
      influenceFor: function (p) {
        let self = this
        let bosses = self.message.bosses[p.id]
        let values = _.values(self.message.playedChips[p.id])
        let count = _.reduce(values, function(sum, n) {
          return sum + n
        }, 0)
        return bosses + count
      }
    },
    computed: {
      pids: function () {
        return _.map(this.message.bosses, function(value, key) {
          return _.toInteger(key)
        })
      },
      ps: function () {
        return this.playersByPIDS(this.pids)
      },
    }
  }
</script>
