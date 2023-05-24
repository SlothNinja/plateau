package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

func (cl Client) finishTurnHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.Current(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	var g game
	if err := sn.GetGame(ctx, cl.Client, &g, cu); err != nil {
		sn.JErr(ctx, err)
		return
	}

	var gc game
	if err := sn.GetCommitted(ctx, cl.Client, &gc); err != nil {
		sn.JErr(ctx, err)
		return
	}

	if gc.Undo.Committed != g.Undo.Committed {
		sn.JErr(ctx, fmt.Errorf("invalid commit: %w", sn.ErrValidation))
		return
	}

	var cp, np *player

	switch g.Phase {
	case bidPhase:
		cp, np, err = g.bidFinishTurn(cu)
	case exchangePhase:
		cp, np, err = g.exchangeFinishTurn(cu)
	case incObjectivePhase:
		cp, np, err = g.incObjectiveFinishTurn(cu)
	case cardPlayPhase:
		cp, np, err = g.playCardFinishTurn(cu)
	default:
		err = fmt.Errorf("cannot finish turn during %q phase: %w", g.Phase, sn.ErrValidation)

	}

	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	cp.Stats.Moves++
	cp.Stats.Think += time.Since(gc.UpdatedAt)

	if np == nil {
		cl.endGame(ctx, g, cu)
		return
	}

	np.reset()
	g.setCurrentPlayers(np)
	err = sn.Commit(ctx, cl.Client, &g, cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	// err = cl.sendTurnNotificationsTo(c, g, g.otherCurrentPlayers(cp)...)
	// err = cl.sendNotifications(c, g)
	// if err != nil {
	// 	cl.Log.Warningf(err.Error())
	// }
	ctx.JSON(http.StatusOK, nil)
}

func (g game) validateFinishTurn(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateCurrentPlayer(cu)
	sn.Debugf("cp: %#v", cp)
	switch {
	case err != nil:
		return nil, err
	case !cp.PerformedAction:
		return nil, fmt.Errorf("%s has yet to perform an action: %w", cu.Name, sn.ErrValidation)
	default:
		return cp, nil
	}
}
