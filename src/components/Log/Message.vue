<template>
  <component
    :is='component'
    :message='message'
    :game='game'
    :pid='pid'
  >
  </component>
</template>

<script>
  import AnnounceWinners from '@/components/log/AnnounceWinners'
  import StartGame from '@/components/log/StartGame'
  import ElectionResults from '@/components/log/ElectionResults'
  import CastleGarden from '@/components/log/CastleGarden'
  import PlaceBosses from '@/components/log/PlaceBosses'
  import AssignOffices from '@/components/log/AssignOffices'
  import AwardChips from '@/components/log/AwardChips'
  import EndGameVP from '@/components/log/EndGameVP'
  import ScoreVP from '@/components/log/ScoreVP'
  import PlaceImmigrant from '@/components/log/PlaceImmigrant'
  import ElectionPlaceImmigrant from '@/components/log/ElectionPlaceImmigrant'
  import RemoveImmigrant from '@/components/log/RemoveImmigrant'
  import MoveImmigrant from '@/components/log/MoveImmigrant'
  import Slander from '@/components/log/Slander'
  import TakeChip from '@/components/log/TakeChip'
  import DeputyTakeChip from '@/components/log/DeputyTakeChip'
  import LockUp from '@/components/log/LockUp'
  import ElectionTakeChip from '@/components/log/ElectionTakeChip'

  export default {
    name: 'sn-log-message',
    props: [ 'message', 'game', 'pid' ],
    components: {
      'sn-log-announce-winners-msg': AnnounceWinners,
      'sn-log-start-game-msg': StartGame,
      'sn-log-election-results-msg': ElectionResults,
      'sn-log-castle-garden-msg': CastleGarden,
      'sn-log-place-bosses-msg': PlaceBosses,
      'sn-log-assign-offices-msg': AssignOffices,
      'sn-log-award-chips-msg': AwardChips,
      'sn-log-end-game-vp-msg': EndGameVP,
      'sn-log-score-vp-msg': ScoreVP,
      'sn-log-place-immigrant-msg': PlaceImmigrant,
      'sn-log-election-place-immigrant-msg': ElectionPlaceImmigrant,
      'sn-log-remove-immigrant-msg': RemoveImmigrant,
      'sn-log-move-immigrant-msg': MoveImmigrant,
      'sn-log-slander-msg': Slander,
      'sn-log-lock-up-msg': LockUp,
      'sn-log-take-chip-msg': TakeChip,
      'sn-log-deputy-take-chip-msg': DeputyTakeChip,
      'sn-log-election-take-chip-msg': ElectionTakeChip,
    },
    computed: {
      component: function () {
        var self = this
        return `sn-log-${self.message.template}-msg`
      }
    }
  }
</script>
