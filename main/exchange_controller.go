package main

import (
	"fmt"
	"net/http"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startExchange() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = exchangePhase

	// Declarer exchanges cards with talon.
	declarer := g.declarer()
	declarer.hand = append(declarer.hand, g.deck...)
	g.deck = nil
	return declarer
}

func (cl Client) exchangeHandler(c *gin.Context) {
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

	err = g.exchange(c, cu)
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

func (g *game) exchange(c *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, cards, err := g.validateExchange(c, cu)
	if err != nil {
		return err
	}

	g.deck = cards
	// remove cards from hand
	_, cp.hand = pie.Diff(cp.hand, cards)
	cp.performedAction = true

	g.newEntryFor(cp.id, message{"template": "card-exchange"})

	g.Undo.Update()
	return nil
}

func (g game) validateExchange(c *gin.Context, cu sn.User) (*player, []card, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePlayerAction(cu)
	if err != nil {
		return nil, nil, err
	}

	cards, err := getCards(c)
	if err != nil {
		return nil, nil, err
	}
	sn.Debugf("cards: %#v", cards)
	sn.Debugf("cp.hand: %#v", cp.hand)

	switch {
	case g.Phase != exchangePhase:
		return nil, nil, fmt.Errorf("cannot exchange cards in %q phase: %w", g.Phase, sn.ErrValidation)
	case g.lastBid().exchange != exchangeBid:
		return nil, nil, fmt.Errorf("winning bid did not include card exchange: %w", sn.ErrValidation)
	case g.NumPlayers == 2 && len(cards) != 2:
		return nil, nil, fmt.Errorf("must exchange two cards: %w", sn.ErrValidation)
	case g.NumPlayers >= 3 && g.NumPlayers <= 6 && len(cards) != 3:
		return nil, nil, fmt.Errorf("must exchange three cards: %w", sn.ErrValidation)
	case !pie.All(cards, func(c card) bool { return pie.Contains(cp.hand, c) }):
		return nil, nil, fmt.Errorf("must exchange from your hand: %w", sn.ErrValidation)
	default:
		return cp, cards, nil
	}
}

func (g *game) exchangeFinishTurn(cu sn.User) (*player, *player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validateExchangeFinishTurn(cu)
	if err != nil {
		return nil, nil, err
	}

	np := g.startPickPartner()
	return cp, np, nil
}

func (g game) validateExchangeFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != exchangePhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", exchangePhase, g.Phase, sn.ErrValidation)
	case g.NumPlayers >= 3 && g.NumPlayers <= 6 && len(cp.hand) != 13:
		return nil, fmt.Errorf("you must have thirteen cards after the exchange: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}
