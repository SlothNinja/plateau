package main

import (
	"fmt"
	"net/http"

	"github.com/SlothNinja/sn/v2"
	"github.com/gin-gonic/gin"
)

func (cl *Client) bidHandler(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{
		"game": g,
		"cu":   cu,
	})
}

func (g *game) placeBid(c *gin.Context, cu *sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, bid, err := g.validatePlaceBid(c, cu)
	if err != nil {
		return err
	}

	sn.Debugf("placeBid bid: %#v", bid)

	cp.performedAction = true
	g.bids = append(g.bids, bid)
	g.newEntryFor(cp.id, message{
		"template": "placed-bid",
		"bid":      bid,
	})
	return nil
}

func (g *game) validatePlaceBid(c *gin.Context, cu *sn.User) (*player, bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// define noBid here, as bid type shadowed by bid variable after getBid call
	noBid := bid{}

	cp, err := g.validatePlayerAction(cu)
	if err != nil {
		return nil, noBid, err
	}

	bid, err := getBid(c)
	if err != nil {
		return nil, noBid, err
	}
	sn.Debugf("getBid bid: %#v", bid)

	bidValue, err := bid.value(g.NumPlayers)
	if err != nil {
		return nil, noBid, err
	}

	minValue, err := minBid(g.NumPlayers).value(g.NumPlayers)
	if err != nil {
		return nil, noBid, err
	}

	currentBid := g.currentBid()
	currentBidValue, err := g.currentBidValue(g.NumPlayers)
	if err != nil {
		return nil, noBid, err
	}

	switch {
	case bidValue < minValue:
		return nil, noBid, fmt.Errorf("bid has value of %d, which is less than the minimum bid of %d: %w",
			bidValue, minValue, sn.ErrValidation)
	case bidValue < currentBidValue:
		return nil, noBid, fmt.Errorf("bid has value of %d, which is less than the current bid of %d: %w",
			bidValue, currentBidValue, sn.ErrValidation)
		// the following should never happen
	case bid.pid == currentBid.pid:
		return nil, noBid, fmt.Errorf("you already have the current bid: %w", sn.ErrValidation)
	default:
		return cp, bid, nil
	}
}
