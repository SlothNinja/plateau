package client

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startIncObjective() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Header.Phase = incObjectivePhase
	pie.Each(g.Players, (*player).bidReset)
	return g.selectIncrementer()
}

func (g *game) selectIncrementer() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return g.NextPlayer(g.declarer(), func(p *player) bool { return pie.Contains(g.declarers(), p) && !p.Bid })
}

func (g *game) partnerPIDS() []sn.PID {
	if len(g.State.DeclarersTeam) < 2 {
		return nil
	}
	return g.State.DeclarersTeam[1:]
}

// func incObjectiveAction(sngame *sn.Game[state, player, *player], ctx *gin.Context, cu sn.User) error {
// 	sn.Debugf(msgEnter)
// 	defer sn.Debugf(msgExit)
//
// 	return (&game{sngame}).incObjective(ctx, cu)
// }

func (g *game) incObjective(ctx *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, bid, err := g.validateIncObjective(ctx, cu)
	if err != nil {
		return err
	}

	cp.PerformedAction = true
	cp.Bid = true
	g.State.Bids = append(g.State.Bids, bid)
	g.NewEntry(incObjectiveTemplate, sn.Entry{"PID": cp.id(), "HandNumber": g.currentHand()}, sn.Line{"Bid": bid})
	return nil
}

const incObjectiveTemplate = "inc-objective"

func (g *game) validateIncObjective(ctx *gin.Context, cu sn.User) (*player, bid, error) {
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

	objValue1 := bid.Objective.value()

	objValue2 := g.lastBid().Objective.value()

	switch {
	case g.Header.Phase != incObjectivePhase:
		return nil, noBid, fmt.Errorf("expected %q phase but have %q phase: %w", incObjectivePhase, g.Header.Phase, sn.ErrValidation)
	case bid.Exchange != g.lastBid().Exchange:
		return nil, noBid, fmt.Errorf("you cannot change the exchange characteristic of the bid: %w", sn.ErrValidation)
	case bid.Teams != g.lastBid().Teams:
		return nil, noBid, fmt.Errorf("you cannot change the teams characteristic of the bid: %w", sn.ErrValidation)
	case bid.PID != cp.id():
		return nil, noBid, fmt.Errorf("you cannot change the declarer of the bid: %w", sn.ErrValidation)
	case objValue1 < objValue2:
		return nil, noBid, fmt.Errorf("you cannot decrease the objective of the bid: %w", sn.ErrValidation)
	default:
		return cp, bid, nil
	}
}

// func incObjectiveFinishTurnAction(sngame *sn.Game[state, player, *player], ctx *gin.Context, cu sn.User) (*player, *player, error) {
// 	sn.Debugf(msgEnter)
// 	defer sn.Debugf(msgExit)
//
// 	return (&game{sngame}).incObjectiveFinishTurn(ctx, cu)
// }

func (g *game) incObjectiveFinishTurn(_ *gin.Context, cu sn.User) (sn.PID, sn.PID, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateIncObjectiveFinishTurn(cu)
	if err != nil {
		return sn.NoPID, sn.NoPID, err
	}

	np := g.selectIncrementer()
	if np == nil {
		np = g.startCardPlay()
	}
	return cp.id(), np.id(), nil
}

func (g *game) validateIncObjectiveFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Header.Phase != incObjectivePhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", incObjectivePhase, g.Header.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}
}
