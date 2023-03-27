package main

import (
	"encoding/json"
	"fmt"

	"github.com/SlothNinja/sn/v2"
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

func (b bid) value(numPlayers int) (int, error) {

	ev, err := b.exchange.value(numPlayers)
	if err != nil {
		return 0, err
	}
	ov, err := b.objective.value()
	if err != nil {
		return 0, err
	}
	tv, err := b.teams.value(numPlayers)
	if err != nil {
		return 0, err
	}
	return ev + ov + tv, nil
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

func (g *game) currentBid() bid {
	return pie.Last(g.bids)
}

func (g *game) currentBidValue(numPlayers int) (int, error) {
	if len(g.bids) == 0 {
		return 0, nil
	}
	return g.currentBid().value(numPlayers)
}

func (e exchange) value(numPlayers int) (int, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)
	sn.Debugf("exchange: %q numPlayers: %d", e, numPlayers)
	switch numPlayers {
	case 2, 3, 4, 5:
		switch e {
		case noExchangeBid:
			return 1, nil
		case exchangeBid:
			return 2, nil
		default:
			return 0, fmt.Errorf("%s is an invalid exchange bid: %w", e, sn.ErrValidation)
		}
	case 6:
		switch e {
		case nullExchangeBid:
			return 0, nil
		default:
			return 0, fmt.Errorf("%s is an invalid exchange bid: %w", e, sn.ErrValidation)
		}
	default:
		return 0, fmt.Errorf("%s is an invalid exchange bid: %w", e, sn.ErrValidation)
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
	bridgeBid    objective = "bridge"
	yBid         objective = "y"
	forkBid      objective = "fork"
	fiveSidesBid objective = "five sides"
	sixSidesBid  objective = "six sides"
)

func (o objective) value() (int, error) {
	switch o {

	case "bridge":
		return 0, nil
	case "y":
		return 2, nil
	case "fork":
		return 4, nil
	case "five sides":
		return 6, nil
	case "six sides":
		return 8, nil
	default:
		return 0, fmt.Errorf("invalid objective bid of: %s: %w", o, sn.ErrValidation)
	}
}

type teams string

const (
	noTeamBid teams = ""
	soloBid   teams = "solo"
	duoBid    teams = "duo"
	trioBid   teams = "trio"
)

func (t teams) value(numPlayers int) (int, error) {
	switch numPlayers {
	case 2, 3:
		switch t {
		case noTeamBid:
			return 0, nil
		default:
			return 0, fmt.Errorf("for %d players, %s is invalid teams bid: %w", numPlayers, t, sn.ErrValidation)
		}
	case 4, 5:
		switch t {
		case duoBid:
			return 0, nil
		case soloBid:
			return 5, nil
		default:
			return 0, fmt.Errorf("for %d players, %s is invalid teams bid: %w", numPlayers, t, sn.ErrValidation)
		}
	case 6:
		switch t {
		case trioBid:
			return 0, nil
		case duoBid:
			return 5, nil
		case soloBid:
			return 10, nil
		default:
			return 0, fmt.Errorf("for %d players, %s is invalid teams bid: %w", numPlayers, t, sn.ErrValidation)
		}
	default:
		return 0, fmt.Errorf("for %d players, %s is invalid teams bid: %w", numPlayers, t, sn.ErrValidation)
	}
}
