package main

import (
	"fmt"
	"net/http"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (cl Client) bidHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(ctx, cu, noUndo)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = g.placeBid(ctx, cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.putCached(ctx, g, g.Undo.Current, cu.ID())
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

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

	// g.newEntryFor(cp.ID, message{
	// 	"Template": "placed-bid",
	// 	"Bid":      bid,
	// })

	g.Undo.Update()
	return nil
}

func (g game) validatePlaceBid(ctx *gin.Context, cu sn.User) (*player, bid, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// define noBid here, as bid type shadowed by bid variable after getBid call
	noBid := bid{}

	cp, err := g.validatePlayerAction(cu)
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

func (cl Client) passBidHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(ctx, cu, noUndo)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = g.passBid(cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.putCached(ctx, g, g.Undo.Current, cu.ID())
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (g *game) passBid(cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePassBid(cu)
	if err != nil {
		return err
	}

	cp.PerformedAction = true
	cp.Bid = true
	cp.Passed = true

	// g.newEntryFor(cp.ID, message{"template": "pass-bid"})

	g.Undo.Update()
	return nil
}

func (g game) validatePassBid(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePlayerAction(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != bidPhase && g.Phase != incObjectivePhase:
		return nil, fmt.Errorf("cannot pass during %q phase: %w", g.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}

}

func (g *game) bidFinishTurn(cu sn.User) (*player, *player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateBidFinishTurn(cu)
	if err != nil {
		return nil, nil, err
	}

	np := g.nextPlayer(cp, func(p *player) bool {
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

func (cl Client) abdicateHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(ctx, cu, noUndo)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = g.abdicate(cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.putCached(ctx, g, g.Undo.Current, cu.ID())
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (g *game) abdicate(cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateAbdicate(cu)
	if err != nil {
		return err
	}

	cp.PerformedAction = true
	cp.Bid = true
	g.DeclarersTeam[0], g.DeclarersTeam[1] = g.DeclarersTeam[1], g.DeclarersTeam[0]

	g.Undo.Update()
	return nil
}

func (g game) validateAbdicate(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePlayerAction(cu)
	switch {
	case err != nil:
		return nil, err
	case g.NumPlayers != 5:
		return nil, fmt.Errorf("the declarer can only abdicate in a 5 player game: %w", sn.ErrValidation)
	case !(pie.Contains(g.DeclarersTeam, g.lastBid().PID) && g.lastBid().PID != cp.ID):
		return nil, fmt.Errorf("the declarer can only abdicate if partner increased objective: %w", sn.ErrValidation)
	case len(g.DeclarersTeam) != 2:
		return nil, fmt.Errorf("the declarer can only abdicate if there is a partner: %w", sn.ErrValidation)
	case pie.First(g.DeclarersTeam) != cp.ID:
		return nil, fmt.Errorf("only the declarer may abdicate: %w", sn.ErrValidation)
	case g.Phase != incObjectivePhase:
		return nil, fmt.Errorf("cannot abdicate during %q phase: %w", g.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}
}
