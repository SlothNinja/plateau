package main

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

func (g *game) placeBid(ctx *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, bid, err := g.validatePlaceBid(ctx, cu)
	if err != nil {
		return err
	}

	cp.PerformedAction = true
	cp.Bid = true
	g.Bids = append(g.Bids, bid)
	g.DeclarersTeam = []sn.PID{cp.ID}

	g.NewEntry(placedBidTemplate, sn.Entry{"PID": cp.ID, "HandNumber": g.currentHand()}, sn.Line{"Bid": bid})

	return nil
}

const placedBidTemplate = "placed-bid"

func (g game) validatePlaceBid(ctx *gin.Context, cu sn.User) (*player, bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// define noBid here, as bid type shadowed by bid variable after getBid call
	noBid := bid{}

	cp, err := g.ValidatePlayerAction(cu)
	if err != nil {
		return nil, noBid, err
	}

	bid, err := g.validateBid(ctx)
	if err != nil {
		return nil, noBid, err
	}

	bidValue := bid.value(g.NumPlayers)

	currentBidValue := g.currentBidValue()

	switch {
	case g.Phase != bidPhase:
		return nil, noBid, fmt.Errorf("expected %q phase but have %q phase: %w", bidPhase, g.Phase, sn.ErrValidation)
	case bidValue <= currentBidValue:
		return nil, noBid, fmt.Errorf("bid has value of %d, which is not greater than the current bid of %d: %w",
			bidValue, currentBidValue, sn.ErrValidation)
	default:
		return cp, bid, nil
	}
}

func (g game) validateBid(ctx *gin.Context) (bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// define noBid here, as bid type shadowed by bid variable after getBid call
	noBid := bid{}

	bid, err := getBid(ctx)
	if err != nil {
		return noBid, err
	}

	bidValue := bid.value(g.NumPlayers)

	minValue := minBid(g.NumPlayers).value(g.NumPlayers)

	if bidValue < minValue {
		return noBid, fmt.Errorf("bid has value of %d, which is less than the minimum bid of %d: %w",
			bidValue, minValue, sn.ErrValidation)
	}
	return bid, nil
}
