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

	cu, err := cl.User.Current(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(ctx, cu, noUndo)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	gc, err := cl.getCommitted(ctx)
	if err != nil {
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
	err = cl.commit(ctx, g, cu)
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
	cp, err := g.validateCurrentPlayer(cu)
	switch {
	case err != nil:
		return nil, err
	case !cp.PerformedAction:
		return nil, fmt.Errorf("%s has yet to perform an action: %w", cu.Name, sn.ErrValidation)
	default:
		return cp, nil
	}
}
