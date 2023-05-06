// Package provides Tammany Hall service
package main

import (
	"errors"
	"fmt"

	"github.com/SlothNinja/sn/v3"
)

func (g game) validatePlayerAction(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateCurrentPlayer(cu)
	switch {
	case err != nil:
		return nil, err
	case cp.PerformedAction:
		return nil, fmt.Errorf("current player already performed action: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}

func (g game) validateCurrentPlayer(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp := g.currentPlayerFor(cu)
	if cp == nil {
		return nil, sn.ErrPlayerNotFound
	}
	return cp, nil
}

func validateAdmin(cu sn.User) error {
	if cu.IsAdmin() {
		return nil
	}
	return errors.New("not admin")
}
