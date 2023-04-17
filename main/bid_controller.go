package main

import (
	"fmt"
	"net/http"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

func (cl Client) bidHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(c, cu, noUndo)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	err = g.placeBid(c, cu)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	err = cl.putCached(c, g, g.Undo.Current, cu.ID())
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": g})
}

func (g *game) placeBid(c *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, bid, err := g.validatePlaceBid(c, cu)
	if err != nil {
		return err
	}

	cp.performedAction = true
	cp.bid = true
	g.bids = append(g.bids, bid)
	g.declarersTeam = []sn.PID{cp.id}

	g.newEntryFor(cp.id, message{
		"template": "placed-bid",
		"bid":      bid,
	})

	g.Undo.Update()
	return nil
}

func (g game) validatePlaceBid(c *gin.Context, cu sn.User) (*player, bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// define noBid here, as bid type shadowed by bid variable after getBid call
	noBid := bid{}

	cp, err := g.validatePlayerAction(cu)
	if err != nil {
		return nil, noBid, err
	}

	bid, err := g.validateBid(c)
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

func (g game) validateBid(c *gin.Context) (bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// define noBid here, as bid type shadowed by bid variable after getBid call
	noBid := bid{}

	bid, err := getBid(c)
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

func (cl Client) passBidHandler(c *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(c)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(c, cu, noUndo)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	err = g.passBid(cu)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	err = cl.putCached(c, g, g.Undo.Current, cu.ID())
	if err != nil {
		sn.JErr(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"game": g})
}

func (g *game) passBid(cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePassBid(cu)
	if err != nil {
		return err
	}

	cp.performedAction = true
	cp.bid = true
	cp.passed = true

	g.newEntryFor(cp.id, message{"template": "pass-bid"})

	g.Undo.Update()
	return nil
}

func (g game) validatePassBid(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return g.validatePlayerAction(cu)
}

func (g *game) bidFinishTurn(cu sn.User) (*player, *player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateBidFinishTurn(cu)
	if err != nil {
		return nil, nil, err
	}

	np := g.nextPlayer(cp, func(p *player) bool {
		return !p.passed && p.id != g.lastBid().pid
	})

	if np.id != sn.NoPID {
		// Proceed to next bidder
		return cp, np, nil
	}

	// Log winning bid
	lastBid := g.lastBid()
	g.newEntryFor(lastBid.pid, message{
		"template": "won-bid",
		"bid":      lastBid,
	})

	// Proceed to next phase
	if lastBid.exchange == exchangeBid {
		np = g.startExchange()
		return cp, np, nil
	}

	if lastBid.includesPartner() {
		np = g.startPickPartner()
		return cp, np, nil
	}

	g.startIncObjective()
	np = g.selectIncrementer(nil)
	return cp, np, nil

}

func (g game) validateBidFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != bidPhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", bidPhase, g.Phase, sn.ErrValidation)
	case !cp.bid:
		return nil, fmt.Errorf("you must bid or pass before finishing turn: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}
