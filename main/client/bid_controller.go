package client

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

const startBidTemplate = "start-bid"

func (g *game) startBidPhase() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Header.Phase = bidPhase
	pie.Each(g.Players, (*player).bidReset)
	g.State.Bids = nil
	g.NewEntry(startBidTemplate, sn.H{
		"PID":        g.forehand().id(),
		"HandNumber": g.currentHand(),
	})
	return g.forehand()
}

func (g *game) bidFinishTurn(ctx *gin.Context, cu sn.User) (sn.PID, sn.PID, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateBidFinishTurn(cu)
	if err != nil {
		return sn.NoPID, sn.NoPID, err
	}

	np := g.NextPlayer(cp, func(p *player) bool {
		return !p.Passed && (p.id() != g.lastBid().PID)
	})

	if np != nil {
		// Proceed to next bidder
		return cp.id(), np.id(), nil
	}

	// all players passed, then next dealer deals new hand
	if pie.All(g.Players, func(p *player) bool { return p.Passed }) {
		np = g.startEndHandPhase(dPush, nil)
		return cp.id(), np.id(), nil
	}

	np = g.startExchange()
	if np != nil {
		return cp.id(), np.id(), nil
	}

	np = g.startPickPartner()
	if np != nil {
		return cp.id(), np.id(), nil
	}

	np = g.startIncObjective()
	return cp.id(), np.id(), nil
}

func (g *game) validateBidFinishTurn(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Header.Phase != bidPhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", bidPhase, g.Header.Phase, sn.ErrValidation)
	case !cp.Bid:
		return nil, fmt.Errorf("you must bid or pass before finishing turn: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}
