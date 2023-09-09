package client

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startBidPhase() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = bidPhase
	pie.Each(g.Players, (*player).bidReset)
	g.Bids = nil
	return g.forehand()
}

func (g *game) bidFinishTurn(ctx *gin.Context, cu sn.User) (*player, *player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateBidFinishTurn(cu)
	if err != nil {
		return nil, nil, err
	}

	np := g.NextPlayer(cp, func(p *player) bool {
		return !p.Passed && p.ID != g.lastBid().PID
	})

	if np != nil {
		// Proceed to next bidder
		return cp, np, nil
	}

	// all players passed, then next dealer deals new hand
	if pie.All(g.Players, func(p *player) bool { return p.Passed }) {
		np = g.startEndHandPhase(dPush, nil)
		return cp, np, nil
	}

	// Log winning bid
	// g.newEntryFor(lastBid.PID, message{
	// 	"template": "won-bid",
	// 	"bid":      lastBid,
	// })

	np = g.startExchange()
	if np != nil {
		return cp, np, nil
	}

	np = g.startPickPartner()
	if np != nil {
		return cp, np, nil
	}

	np = g.startIncObjective()
	return cp, np, nil
}

func (g game) validateBidFinishTurn(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != bidPhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", bidPhase, g.Phase, sn.ErrValidation)
	case !cp.Bid:
		return nil, fmt.Errorf("you must bid or pass before finishing turn: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}
