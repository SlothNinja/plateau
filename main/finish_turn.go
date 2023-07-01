package main

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
)

func (g game) validateFinishTurn(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.ValidateCurrentPlayer(cu)
	switch {
	case err != nil:
		return nil, err
	case !cp.PerformedAction:
		return nil, fmt.Errorf("%s has yet to perform an action: %w", cu.Name, sn.ErrValidation)
	default:
		return cp, nil
	}
}
