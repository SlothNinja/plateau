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
	g.DeclarersTeam[0], g.DeclarersTeam[1] = g.DeclarersTeam[1], g.DeclarersTeam[0]

	return nil
}

func (g game) validateAbdicate(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.ValidatePlayerAction(cu)
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
