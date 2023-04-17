package main

import (
	"fmt"
	"net/http"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startCardPlay() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = cardPlayPhase
	return g.forehand()
}

func (cl Client) playCardHandler(c *gin.Context) {
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

	err = g.playCard(c, cu)
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

func (g *game) playCard(c *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// warning card variable shadows card type
	cp, card, err := g.validatePlayCard(c, cu)
	if err != nil {
		return err
	}

	// need to remove card from hand, before updating card.playedBy
	// otherise, card will not match card in hand and therefore will not be removed
	cp.hand = removeCards(cp.hand, card)

	card.playedBy = cp.id
	g.tricks[g.trickNumber()].cards = append(g.currentTrick().cards, card)

	cp.performedAction = true

	g.newEntryFor(cp.id, message{"template": "card-exchange"})

	g.Undo.Update()
	return nil
}

func removeCards(cards []card, remove ...card) []card {
	_, remainingCards := pie.Diff(cards, remove)
	return remainingCards
}

func (g game) validatePlayCard(c *gin.Context, cu sn.User) (*player, card, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	noCard := card{}
	cp, err := g.validatePlayerAction(cu)
	if err != nil {
		return nil, noCard, err
	}

	cards, err := getCards(c)
	if err != nil {
		return nil, noCard, err
	}

	if len(cards) != 1 {
		return nil, noCard, fmt.Errorf("must play one and only one card: %w", sn.ErrValidation)
	}

	playedCard := pie.First(cards)
	ledSuit := g.ledSuit()

	switch {
	case g.Phase != cardPlayPhase:
		return nil, noCard, fmt.Errorf("cannot play cards in %q phase: %w", g.Phase, sn.ErrValidation)
	case !cp.hasCards(playedCard):
		return nil, noCard, fmt.Errorf("must play card from your hand: %w", sn.ErrValidation)
	case ledSuit != noSuit && cp.hasSuit(ledSuit) && playedCard.suit != ledSuit:
		return nil, noCard, fmt.Errorf("must play card of %q, which is led suit: %w", ledSuit, sn.ErrValidation)
	case ledSuit != noSuit && !cp.hasSuit(ledSuit) && cp.hasSuit(trumps) && playedCard.suit != trumps:
		return nil, noCard, fmt.Errorf("must play trump card if you do not have a card of led suit: %w", sn.ErrValidation)
	default:
		return cp, playedCard, nil
	}
}

func (p player) hasCards(cards ...card) bool {
	return pie.All(cards, func(c card) bool { return pie.Contains(p.hand, c) })
}

func (p player) hasSuit(s suit) bool {
	return pie.Any(p.hand, func(c card) bool { return c.suit == s })
}

func (g game) ledSuit() suit {
	return pie.First(g.currentTrick().cards).suit
}

func (g *game) playCardFinishTurn(cu sn.User) (*player, *player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePlayCardFinishTurn(cu)
	if err != nil {
		return nil, nil, err
	}

	var np *player
	if len(g.currentTrick().cards) != g.NumPlayers {
		np = g.nextPlayer(cp)
		return cp, np, nil
	}

	np = g.endTrick()
	endHand, success, _ := g.endHandWithReveal()
	if !endHand {
		return cp, np, nil
	}

	g.startEndHandPhase()
	g.scoreHand(success)

	sn.Warningf("Still need to check for end game")
	endGame := false
	if !endGame {
		g.startHand()
		np = g.startBidPhase()
	}
	return cp, np, nil
}

func (g game) validatePlayCardFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Phase != cardPlayPhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", cardPlayPhase, g.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}
}
