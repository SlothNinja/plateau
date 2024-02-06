package client

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) abdicate(_ *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateAbdicate(cu)
	if err != nil {
		return err
	}

	cp.PerformedAction = true
	cp.Bid = true
	g.State.DeclarersTeam[0], g.State.DeclarersTeam[1] = g.State.DeclarersTeam[1], g.State.DeclarersTeam[0]

	return nil
}

func (g *game) validateAbdicate(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.ValidatePlayerAction(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Header.NumPlayers != 5:
		return nil, fmt.Errorf("the declarer can only abdicate in a 5 player game: %w", sn.ErrValidation)
	case !(pie.Contains(g.State.DeclarersTeam, g.lastBid().PID) && g.lastBid().PID != cp.id()):
		return nil, fmt.Errorf("the declarer can only abdicate if partner increased objective: %w", sn.ErrValidation)
	case len(g.State.DeclarersTeam) != 2:
		return nil, fmt.Errorf("the declarer can only abdicate if there is a partner: %w", sn.ErrValidation)
	case pie.First(g.State.DeclarersTeam) != cp.id():
		return nil, fmt.Errorf("only the declarer may abdicate: %w", sn.ErrValidation)
	case g.Header.Phase != incObjectivePhase:
		return nil, fmt.Errorf("cannot abdicate during %q phase: %w", g.Header.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}
}
