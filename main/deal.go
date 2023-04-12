package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

func (g *game) deal() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	if g.NumPlayers == 2 {
		g.twoPlayerDeal()
	}
	g.normalDeal()
}

// assumes deck is shuffled
func (g *game) normalDeal() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	const handSize = 13

	// splits deck into chunks of 13 cards, with possible remainder left for deck
	hands := pie.Chunk(g.deck, handSize)

	// assign chunks/hands to players
	for i := range g.players {
		g.players[i].hand = hands[i]
	}

	// if a remainder chunk, assign to deck
	if len(hands) > g.NumPlayers {
		g.deck = hands[g.NumPlayers]
	}
}

// assumes deck is shuffled
func (g *game) twoPlayerDeal() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	sn.Warningf("2P not yet implemented")

}
