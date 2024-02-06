package client

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

	if g.Header.Phase == incObjectivePhase {
		g.NewEntry(passedIncObjectiveTemplate, sn.H{"PID": cp.id(), "HandNumber": g.currentHand()})
	} else {
		g.NewEntry(passedBidTemplate, sn.H{"PID": cp.id(), "HandNumber": g.currentHand()})
	}

	return nil
}

const passedIncObjectiveTemplate = "passed-inc-objective"
const passedBidTemplate = "passed-bid"

func (g *game) validatePassBid(cu sn.User) (*player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.ValidatePlayerAction(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Header.Phase != bidPhase && g.Header.Phase != incObjectivePhase:
		return nil, fmt.Errorf("cannot pass during %q phase: %w", g.Header.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}

}
