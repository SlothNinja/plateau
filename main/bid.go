package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

type bid struct {
	Exchange  exchange
	Objective objective
	Teams     teams
	PID       sn.PID
}

type exchange string

const (
	nullExchangeBid exchange = ""
	noExchangeBid   exchange = "no exchange"
	exchangeBid     exchange = "exchange"
)

func (b bid) value(numPlayers int) int64 {
	return b.Exchange.value(numPlayers) + b.Objective.value() + b.Teams.value(numPlayers)
}

func minBid(numPlayers int) bid {
	switch numPlayers {
	case 2:
		return bid{Exchange: exchangeBid, Objective: yBid}
	case 3:
		return bid{Exchange: exchangeBid, Objective: bridgeBid}
	case 4:
		return bid{Exchange: exchangeBid, Objective: yBid, Teams: duoBid}
	case 5:
		return bid{Exchange: exchangeBid, Objective: bridgeBid, Teams: duoBid}
	case 6:
		return bid{Exchange: noExchangeBid, Objective: yBid, Teams: trioBid}
	default:
		return bid{}
	}
}

func getBid(ctx *gin.Context) (bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var obj bid
	err := ctx.ShouldBind(&obj)
	if err != nil {
		return bid{}, err
	}
	return obj, nil
}

func (g game) currentBid() bid {
	return pie.Last(g.Bids)
}

func (g game) currentBidValue() int64 {
	if len(g.Bids) == 0 {
		return 0
	}
	return g.currentBid().value(g.NumPlayers)
}

func (e exchange) value(numPlayers int) int64 {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	switch numPlayers {
	case 2, 3, 4, 5:
		switch e {
		case noExchangeBid:
			return 2
		case exchangeBid:
			return 1
		default:
			sn.Warningf("%s is an invalid exchange bid", e)
			return 0
		}
	case 6:
		switch e {
		case nullExchangeBid:
			return 0
		default:
			sn.Warningf("%s is an invalid exchange bid", e)
			return 0
		}
	default:
		sn.Warningf("%s is an invalid exchange bid", e)
		return 0
	}
}

type objective string

const (
	noObjectiveBid objective = ""
	bridgeBid      objective = "bridge"
	yBid           objective = "y"
	forkBid        objective = "fork"
	fiveSidesBid   objective = "five sides"
	sixSidesBid    objective = "six sides"
)

func (o objective) value() int64 {
	switch o {
	case "bridge":
		return 0
	case "y":
		return 2
	case "fork":
		return 4
	case "five sides":
		return 6
	case "six sides":
		return 8
	default:
		sn.Warningf("invalid objective bid of: %s", o)
		return 0
	}
}

type teams string

const (
	noTeamBid teams = ""
	soloBid   teams = "solo"
	duoBid    teams = "duo"
	trioBid   teams = "trio"
)

func (t teams) value(numPlayers int) int64 {
	switch numPlayers {
	case 2, 3:
		switch t {
		case noTeamBid:
			return 0
		default:
			sn.Warningf("for %d players, %s is invalid teams bid", numPlayers, t)
			return 0
		}
	case 4, 5:
		switch t {
		case duoBid:
			return 0
		case soloBid:
			return 5
		default:
			sn.Warningf("for %d players, %s is invalid teams bid", numPlayers, t)
			return 0
		}
	case 6:
		switch t {
		case trioBid:
			return 0
		case duoBid:
			return 5
		case soloBid:
			return 10
		default:
			sn.Warningf("for %d players, %s is invalid teams bid", numPlayers, t)
			return 0
		}
	default:
		sn.Warningf("for %d players, %s is invalid teams bid", numPlayers, t)
		return 0
	}
}

func (g game) lastBid() bid {
	return pie.Last(g.Bids)
}

func (b bid) includesPartner() bool {
	return b.Teams == duoBid || b.Teams == trioBid
}
