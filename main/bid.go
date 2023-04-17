package main

import (
	"encoding/json"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

type bid struct {
	exchange  exchange
	objective objective
	teams     teams
	pid       sn.PID
}

type exchange string

const (
	nullExchangeBid exchange = ""
	noExchangeBid   exchange = "no exchange"
	exchangeBid     exchange = "exchange"
)

func (b bid) value(numPlayers int) int {
	return b.exchange.value(numPlayers) + b.objective.value() + b.teams.value(numPlayers)
}

func minBid(numPlayers int) bid {
	switch numPlayers {
	case 2:
		return bid{exchange: exchangeBid, objective: yBid}
	case 3:
		return bid{exchange: exchangeBid, objective: bridgeBid}
	case 4:
		return bid{exchange: exchangeBid, objective: yBid, teams: duoBid}
	case 5:
		return bid{exchange: exchangeBid, objective: bridgeBid, teams: duoBid}
	case 6:
		return bid{exchange: noExchangeBid, objective: yBid, teams: trioBid}
	default:
		return bid{}
	}
}

func getBid(c *gin.Context) (bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	obj := jBid{}
	err := c.ShouldBind(&obj)
	if err != nil {
		return bid{}, err
	}
	sn.Debugf("obj: %#v", obj)
	return bidFrom(obj), nil
}

func (g game) currentBid() bid {
	return pie.Last(g.bids)
}

func (g game) currentBidValue() int {
	if len(g.bids) == 0 {
		return 0
	}
	return g.currentBid().value(g.NumPlayers)
}

func (e exchange) value(numPlayers int) int {
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

type jBid struct {
	Exchange  exchange  `json:"exchange"`
	Objective objective `json:"objective"`
	Teams     teams     `json:"teams"`
	PID       sn.PID    `json:"pid"`
}

func (b bid) MarshalJSON() ([]byte, error) {
	return json.Marshal(jBid{
		Exchange:  b.exchange,
		Objective: b.objective,
		Teams:     b.teams,
		PID:       b.pid,
	})
}

func (b *bid) UnmarshalJSON(bs []byte) error {
	obj := new(jBid)
	err := json.Unmarshal(bs, obj)
	if err != nil {
		return err
	}
	*b = bidFrom(*obj)
	return nil
}

func bidFrom(obj jBid) bid {
	return bid{
		exchange:  obj.Exchange,
		objective: obj.Objective,
		teams:     obj.Teams,
		pid:       obj.PID,
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

func (o objective) value() int {
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

func (t teams) value(numPlayers int) int {
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
	return pie.Last(g.bids)
}

func (b bid) includesPartner() bool {
	return b.teams == duoBid || b.teams == trioBid
}
