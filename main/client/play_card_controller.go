package client

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startCardPlay() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Header.Phase = cardPlayPhase
	return g.forehand()
}

func (g *game) playCard(ctx *gin.Context, cu sn.User) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, c, err := g.validatePlayCard(ctx, cu)
	if err != nil {
		return err
	}

	// need to play/remove card from hand, before updating card.playedBy
	// otherise, card will not match card in hand and therefore will not be removed
	cp.play(c)

	c.PlayedBy = cp.id()
	g.State.Tricks[g.trickIndex()].Cards = append(g.currentTrick().Cards, c)

	cp.PerformedAction = true

	if i := len(g.currentTrick().Cards) - 1; i == 0 {
		g.NewEntry(
			playedCardTemplate,
			sn.H{
				"TrickNumber":   g.trickNumber(),
				"HandNumber":    g.currentHand(),
				"Trick":         g.currentTrick(),
				"DeclarersTeam": g.State.DeclarersTeam},
		)
	} else {
		g.UpdateLastEntry(sn.H{"Trick": g.currentTrick()})
	}

	return nil
}

const playedCardTemplate = "played-card"

func (p *player) play(c card) {
	if pie.Contains(p.Hand, c) {
		p.Hand = removeCards(p.Hand, c)
		return
	}

	if pie.Last(p.Stack0) == c {
		if len(p.Stack0) == 1 {
			p.Stack0 = nil
			return
		}
		p.Stack0 = p.Stack0[0:1]
		return
	}

	if pie.Last(p.Stack1) == c {
		if len(p.Stack1) == 1 {
			p.Stack1 = nil
			return
		}
		p.Stack1 = p.Stack1[0:1]
		return
	}

	if pie.Last(p.Stack2) == c {
		if len(p.Stack2) == 1 {
			p.Stack2 = nil
			return
		}
		p.Stack2 = p.Stack2[0:1]
		return
	}

	if pie.Last(p.Stack3) == c {
		if len(p.Stack3) == 1 {
			p.Stack3 = nil
			return
		}
		p.Stack3 = p.Stack3[0:1]
		return
	}

	if pie.Last(p.Stack4) == c {
		if len(p.Stack4) == 1 {
			p.Stack4 = nil
			return
		}
		p.Stack4 = p.Stack4[0:1]
		return
	}
}

func removeCards(cards []card, remove ...card) []card {
	_, remainingCards := pie.Diff(cards, remove)
	return remainingCards
}

func (g *game) validatePlayCard(ctx *gin.Context, cu sn.User) (*player, card, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	noCard := card{}
	cp, err := g.ValidatePlayerAction(cu)
	if err != nil {
		return nil, noCard, err
	}

	cards, err := getCards(ctx)
	if err != nil {
		return nil, noCard, err
	}

	if len(cards) != 1 {
		return nil, noCard, fmt.Errorf("must play one and only one card: %w", sn.ErrValidation)
	}

	playedCard := pie.First(cards)
	ledSuit := g.ledSuit()

	switch {
	case g.Header.Phase != cardPlayPhase:
		return nil, noCard, fmt.Errorf("cannot play cards in %q phase: %w", g.Header.Phase, sn.ErrValidation)
	case !cp.hasCard(playedCard):
		return nil, noCard, fmt.Errorf("must play card from your hand: %w", sn.ErrValidation)
	case ledSuit != noSuit && cp.hasSuit(ledSuit) && playedCard.Suit != ledSuit:
		return nil, noCard, fmt.Errorf("must play card of %q, which is led suit: %w", ledSuit, sn.ErrValidation)
	case ledSuit != noSuit && !cp.hasSuit(ledSuit) && cp.hasSuit(trumps) && playedCard.Suit != trumps:
		return nil, noCard, fmt.Errorf("must play trump card if you do not have a card of led suit: %w", sn.ErrValidation)
	default:
		return cp, playedCard, nil
	}
}

func (g *game) ledSuit() suit {
	return pie.First(g.currentTrick().Cards).Suit
}

func (g *game) playCardFinishTurn(_ *gin.Context, cu sn.User) (sn.PID, sn.PID, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	cp, err := g.validatePlayCardFinishTurn(cu)
	if err != nil {
		return sn.NoPID, sn.NoPID, err
	}

	var np *player
	if len(g.currentTrick().Cards) != g.Header.NumPlayers {
		np = g.NextPlayer(cp)
		return cp.id(), np.id(), nil
	}

	np = g.endTrick()
	if np != nil {
		g.UpdateLastEntry(sn.H{"WonTrick": np.id()})
	}
	endHand, result, path := g.endHandCheck()
	if !endHand {
		return cp.id(), np.id(), nil
	}

	np = g.startEndHandPhase(result, path)

	return cp.id(), np.id(), nil
}

func (g *game) validatePlayCardFinishTurn(cu sn.User) (*player, error) {
	cp, err := g.validateFinishTurn(cu)
	switch {
	case err != nil:
		return nil, err
	case g.Header.Phase != cardPlayPhase:
		return nil, fmt.Errorf("expected %q phase but have %q phase: %w", cardPlayPhase, g.Header.Phase, sn.ErrValidation)
	default:
		return cp, nil
	}
}
