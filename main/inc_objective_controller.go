package main

import (
	"fmt"
	"net/http"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startIncObjective() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = incObjectivePhase
}

func (g game) selectIncrementer(inc *player) *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	switch len(g.partnerPIDS()) {
	case 0:
		return g.declarer()
	case 1:
		if inc == nil {
			return g.partner()
		}
		return g.declarer()
	case 2:
		switch {
		case inc == nil:
			return g.nextPlayer(g.declarer(), func(p *player) bool { return pie.Contains(g.partners(), p) })
		case inc == g.nextPlayer(g.declarer(), func(p *player) bool { return pie.Contains(g.partners(), p) }):
			return g.nextPlayer(inc, func(p *player) bool { return pie.Contains(g.partners(), p) })
		default:
			return g.declarer()
		}
	default:
		sn.Warningf("len(g.partnerPIDS): %d", len(g.partnerPIDS()))
		return nil
	}
}

func (g game) partnerPIDS() []sn.PID {
	if len(g.declarersTeam) < 2 {
		return nil
	}
	return g.declarersTeam[1:]
}

func (cl Client) incObjectiveHandler(c *gin.Context) {
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

	err = g.incObjective(c, cu)
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

func (g *game) incObjective(c *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, bid, err := g.validateIncObjective(c, cu)
	if err != nil {
		return err
	}

	cp.performedAction = true
	if g.lastBid().objective == bid.objective {
		g.newEntryFor(cp.id, message{
			"template": "no-increased-objective",
		})
	}

	g.bids = append(g.bids, bid)
	g.newEntryFor(cp.id, message{
		"template": "increased-objective",
		"bid":      bid,
	})

	g.Undo.Update()
	return nil
}

func (g game) validateIncObjective(c *gin.Context, cu sn.User) (*player, bid, error) {
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

	objValue1 := bid.objective.value()

	objValue2 := g.lastBid().objective.value()

	switch {
	case g.Phase != incObjectivePhase:
		return nil, noBid, fmt.Errorf("expected %q phase but have %q phase: %w", incObjectivePhase, g.Phase, sn.ErrValidation)
	case bid.exchange != g.lastBid().exchange:
		return nil, noBid, fmt.Errorf("you cannot change the exchange characteristic of the bid: %w", sn.ErrValidation)
	case bid.teams != g.lastBid().teams:
		return nil, noBid, fmt.Errorf("you cannot change the teams characteristic of the bid: %w", sn.ErrValidation)
	case bid.pid != g.lastBid().pid:
		return nil, noBid, fmt.Errorf("you cannot change the declarer of the bid: %w", sn.ErrValidation)
	case objValue1 < objValue2:
		return nil, noBid, fmt.Errorf("you cannot decrease the objective of the bid: %w", sn.ErrValidation)
	default:
		return cp, bid, nil
	}
}

func (g *game) incObjectiveFinishTurn(cu sn.User) (*player, *player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateIncObjectiveFinishTurn(cu)
	if err != nil {
		return nil, nil, err
	}

	np := g.selectIncrementer(cp)
	if cp == np {
		np = g.startCardPlay()
	}
	return cp, np, nil
}

func (g game) validateIncObjectiveFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != incObjectivePhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", incObjectivePhase, g.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}
}
