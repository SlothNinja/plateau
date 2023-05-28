package main

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/gin-gonic/gin"
)

func (g *game) passBid(_ *gin.Context, cu sn.User) error {
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

	return nil
}

func (g game) validatePassBid(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.ValidatePlayerAction(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != bidPhase && g.Phase != incObjectivePhase:
		return nil, fmt.Errorf("cannot pass during %q phase: %w", g.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}

}
