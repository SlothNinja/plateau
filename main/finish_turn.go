package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

func (cl *Client) finishTurnHandler(c *gin.Context) {
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

	gc, err := cl.getCommitted(c)
	if err != nil {
		sn.JErr(c, err)
		return
	}

	if gc.Undo.Committed != g.Undo.Committed {
		sn.JErr(c, fmt.Errorf("invalid commit: %w", sn.ErrValidation))
		return
	}

	var cp, np *player

	switch g.Phase {
	case bidPhase:
		cp, np, err = g.bidFinishTurn(c, cu)
	case exchangePhase:
		cp, np, err = g.exchangeFinishTurn(c, cu)
	case incObjectivePhase:
		cp, np, err = g.incObjectiveFinishTurn(c, cu)
	case cardPlayPhase:
		cp, np, err = g.playCardFinishTurn(c, cu)
	default:
		err = fmt.Errorf("cannot finish turn during %q phase: %w", g.Phase, sn.ErrValidation)

	}

	if err != nil {
		sn.JErr(c, err)
		return
	}

	cp.stats.Moves++
	cp.stats.Think += time.Since(gc.UpdatedAt)

	if np != nil {
		np.reset()
		// g.beginningOfTurnReset()
		g.setCurrentPlayers(np)
	}

	err = cl.commit(c, g, cu.ID())
	if err != nil {
		sn.JErr(c, err)
		return
	}

	// err = cl.sendTurnNotificationsTo(c, g, g.otherCurrentPlayers(cp)...)
	err = cl.sendNotifications(c, g)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"game": g})

}

func (g *game) validateFinishTurn(c *gin.Context, cu *sn.User) (*player, error) {
	cp, err := g.validateCurrentPlayer(cu)
	switch {
	case err != nil:
		return nil, err
	case !cp.performedAction:
		return nil, fmt.Errorf("%s has yet to perform an action: %w", cu.Name, sn.ErrValidation)
	default:
		return cp, nil
	}
}
