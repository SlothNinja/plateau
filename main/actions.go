// Package provides Tammany Hall service
package main

import (
	"errors"
	"fmt"

	"github.com/SlothNinja/sn/v3"
)

func (g *game) validatePlayerAction(cu *sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateCurrentPlayer(cu)
	switch {
	case err != nil:
		return nil, err
	case cp.performedAction:
		return nil, fmt.Errorf("current player already performed action: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}

func (g *game) validateCurrentPlayer(cu *sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp := g.currentPlayerFor(cu)
	switch {
	case cu == nil:
		return nil, sn.ErrNotFound
	case cp == nil:
		return nil, sn.ErrPlayerNotFound
	default:
		return cp, nil
	}
}

func (g *game) validateCPorAdmin(cu *sn.User) (*player, error) {
	cp, err := g.validateCurrentPlayer(cu)
	if err == nil {
		return cp, nil
	}

	err = validateAdmin(cu)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func validateAdmin(cu *sn.User) error {
	switch {
	case cu == nil:
		return sn.ErrNotFound
	case !cu.Admin:
		return errors.New("not admin")
	default:
		return nil
	}
}
