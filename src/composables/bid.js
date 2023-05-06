import { unref } from 'vue'
import _get from 'lodash/get'

export function bidValue(game, bid) {
  const g = unref(game)
  const b = unref(bid)
  return exchangeValue(b) + objectiveValue(b) + teamsValue(g, b)
}

export function exchangeValue(bid) {
  const exchange = _get(unref(bid), 'Exchange', '')
  switch (exchange) {
    case 'exchange':
      return 1
    case 'no exchange':
      return 2
    default:
      return 0
  }
}

export function  objectiveValue(bid) {
  switch (_get(unref(bid), 'Objective', '')) {
    case 'bridge':
      return 0
    case 'y':
      return 2
    case 'fork':
      return 4
    case '5-sides':
      return 6
    case '6-sides':
      return 8
    default:
      return 0
  }
}

export function teamsValue(game, bid) {
  switch (_get(unref(game), 'NumPlayers', 0)) {
    case 4:
      return team45(bid)
    case 5:
      return team45(bid)
    case 6:
      return team6(bid)
    default:
      return 0
  }
}

function team45(bid) {
  return (bid.Teams == 'solo') ? 5 : 0
}

function team6(bid) {
  switch (bid.Teams) {
    case 'duo':
      return 5
    case 'solo':
      return 10
    default:
      return 0
  }
}

export function minBid(numPlayers) {
  switch (unref(numPlayers)) {
    case 2:
      return { Exchange: 'exchange', Objective: 'y' }
    case 3:
      return { Exchange: 'exchange', Objective: 'bridge' }
    case 4:
      return { Exchange: 'exchange', Objective: 'y', Teams: 'duo' }
    case 5:
      return { Exchange: 'exchange', Objective: 'bridge', Teams: 'duo' }
    case 6:
      return { Objective: 'y', Teams: 'trio' }
    default:
      return { Exchange: '', Objective: '', Teams: '' }
  }
}
