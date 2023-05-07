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

	if g.lastBid().Exchange == noExchangeBid {
		return nil
	}

	g.Phase = exchangePhase

	// Declarer exchanges cards with talon.
	declarer := g.declarer()
	declarer.Hand = append(declarer.Hand, g.Deck...)
	g.Deck = nil
	return declarer
}

func (cl Client) exchangeHandler(ctx *gin.Context) {
	cl.Log.Debugf(msgEnter)
	defer cl.Log.Debugf(msgExit)

	cu, err := cl.User.Current(ctx)
	if err != nil {
		cl.Log.Warningf(err.Error())
	}

	g, err := cl.getGame(ctx, cu, noUndo)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = g.exchange(ctx, cu)
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	err = cl.putCached(ctx, g, g.Undo.Current, cu.ID())
	if err != nil {
		sn.JErr(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"game": g})
}

func (g *game) exchange(ctx *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, cards, err := g.validateExchange(ctx, cu)
	if err != nil {
		return err
	}

	g.Deck = cards
	// remove cards from hand
	_, cp.Hand = pie.Diff(cp.Hand, cards)
	cp.PerformedAction = true

	// g.newEntryFor(cp.ID, message{"template": "card-exchange"})

	g.Undo.Update()
	return nil
}

func (g game) validateExchange(ctx *gin.Context, cu sn.User) (*player, []card, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePlayerAction(cu)
	if err != nil {
		return nil, nil, err
	}

	cards, err := getCards(ctx)
	if err != nil {
		return nil, nil, err
	}

	switch {
	case g.Phase != exchangePhase:
		return nil, nil, fmt.Errorf("cannot exchange cards in %q phase: %w", g.Phase, sn.ErrValidation)
	case g.lastBid().Exchange != exchangeBid:
		return nil, nil, fmt.Errorf("winning bid did not include card exchange: %w", sn.ErrValidation)
	case g.NumPlayers == 2 && len(cards) != 2:
		return nil, nil, fmt.Errorf("must exchange two cards: %w", sn.ErrValidation)
	case g.NumPlayers >= 3 && g.NumPlayers <= 6 && len(cards) != 3:
		return nil, nil, fmt.Errorf("must exchange three cards: %w", sn.ErrValidation)
	case !pie.All(cards, func(c card) bool { return pie.Contains(cp.Hand, c) }):
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
	if np != nil {
		return cp, np, nil
	}

	np = g.startIncObjective()
	return cp, np, nil
}

func (g game) validateExchangeFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != exchangePhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", exchangePhase, g.Phase, sn.ErrValidation)
	case g.NumPlayers >= 3 && g.NumPlayers <= 6 && len(cp.Hand) != 13:
		return nil, fmt.Errorf("you must have thirteen cards after the exchange: %w", sn.ErrValidation)
	default:
		return cp, nil
	}
}
