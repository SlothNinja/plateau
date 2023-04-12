package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

func (g *game) endHand() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = endHandPhase
	if g.allCardsPlayed() {
		g.revealTalon()
	}
	g.scoreHand()
	return g.startHand()
}

func (g *game) scoreHand() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	dtbv := g.currentBidValue()
	dl := len(g.declarersTeam)
	ol := g.NumPlayers - dl
	if !g.madeObjective() {
		dtbv = -dtbv
	}

	pie.Each(g.opposers(), func(p *player) { p.score -= dtbv })
	switch {
	case dl == ol:
		pie.Each(g.declarers(), func(p *player) { p.score += dtbv })
	case dl == 1:
		g.declarer().score += dtbv * ol
	case dl == 2 && g.NumPlayers == 5:
		g.declarer().score += dtbv * 2
		g.partner().score += dtbv
	case dl == 2 && g.NumPlayers == 6:
		pie.Each(g.declarers(), func(p *player) { p.score += dtbv * 2 })
	}
}

func (g *game) startHand() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.deck = deckFor(g.NumPlayers)
	g.Phase = dealPhase

	g.resetTrickNumber()

	if g.incHandNumber() == 1 {
		g.randomSeats()
	} else {
		g.newDealer()
	}

	switch g.NumPlayers {
	case 2:
		g.tricks = make([]trick, 17) // 16 trick plus talon
	case 6:
		g.tricks = make([]trick, 13) // 13 trick no talon
	default:
		g.tricks = make([]trick, 14) // 13 trick plus talon
	}

	g.declarersTeam = nil

	g.deal()
	return g.startBidPhase()
}

func (g *game) revealTalon() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var talon trick
	talon.cards = g.deck
	if g.currentBid().exchange == noExchangeBid {
		talon.wonBy = pie.First(g.declarersTeam)
	} else {
		talon.wonBy = pie.First(g.opposersTeam())
	}
	g.tricks = append(g.tricks, talon)
}
