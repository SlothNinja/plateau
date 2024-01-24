package client

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startExchange() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	if g.lastBid().Exchange == noExchangeBid {
		return nil
	}

	g.Header.Phase = exchangePhase

	// Declarer exchanges cards with talon.
	declarer := g.declarer()
	declarer.Hand = append(declarer.Hand, g.State.Deck...)
	g.State.Deck = nil
	return declarer
}

// func exchangeAction(sngame *sn.Game[state, player, *player], ctx *gin.Context, cu sn.User) error {
// 	sn.Debugf(msgEnter)
// 	defer sn.Debugf(msgExit)
//
// 	return (&game{sngame}).exchange(ctx, cu)
// }

func (g *game) exchange(ctx *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, cards, err := g.validateExchange(ctx, cu)
	if err != nil {
		return err
	}

	g.State.Deck = cards
	// remove cards from hand
	_, cp.Hand = pie.Diff(cp.Hand, cards)
	cp.PerformedAction = true

	g.NewEntry(exchangedCardsTemplate, sn.Entry{"PID": cp.id(), "HandNumber": g.currentHand()}, nil)

	return nil
}

const exchangedCardsTemplate = "exchanged-cards"

func (g *game) validateExchange(ctx *gin.Context, cu sn.User) (*player, []card, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.ValidatePlayerAction(cu)
	if err != nil {
		return nil, nil, err
	}

	cards, err := getCards(ctx)
	if err != nil {
		return nil, nil, err
	}

	switch {
	case g.Header.Phase != exchangePhase:
		return nil, nil, fmt.Errorf("cannot exchange cards in %q phase: %w", g.Header.Phase, sn.ErrValidation)
	case g.lastBid().Exchange != exchangeBid:
		return nil, nil, fmt.Errorf("winning bid did not include card exchange: %w", sn.ErrValidation)
	case g.Header.NumPlayers == 2 && len(cards) != 2:
		return nil, nil, fmt.Errorf("must exchange two cards: %w", sn.ErrValidation)
	case g.Header.NumPlayers >= 3 && g.Header.NumPlayers <= 6 && len(cards) != 3:
		return nil, nil, fmt.Errorf("must exchange three cards: %w", sn.ErrValidation)
	case !pie.All(cards, func(c card) bool { return pie.Contains(cp.Hand, c) }):
		return nil, nil, fmt.Errorf("must exchange from your hand: %w", sn.ErrValidation)
	default:
		return cp, cards, nil
	}
}

// func exchangeFinishTurnAction(sngame *sn.Game[state, player, *player], ctx *gin.Context, cu sn.User) (*player, *player, error) {
// 	sn.Debugf(msgEnter)
// 	defer sn.Debugf(msgExit)
//
// 	return (&game{sngame}).exchangeFinishTurn(ctx, cu)
// }

func (g *game) exchangeFinishTurn(_ *gin.Context, cu sn.User) (sn.PID, sn.PID, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateExchangeFinishTurn(cu)
	if err != nil {
		return sn.NoPID, sn.NoPID, err
	}

	np := g.startPickPartner()
	if np != nil {
		return cp.id(), np.id(), nil
	}

	np = g.startIncObjective()
	return cp.id(), np.id(), nil
}

func (g *game) validateExchangeFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Header.Phase != exchangePhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", exchangePhase, g.Header.Phase, sn.ErrValidation)
	case g.Header.NumPlayers >= 3 && g.Header.NumPlayers <= 6 && len(cp.Hand) != 13:
		return nil, fmt.Errorf("you must have thirteen cards after the exchange: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}
