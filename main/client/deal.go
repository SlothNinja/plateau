package client

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

const dealTemplate = "announce-deal"

func (g *game) deal() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	switch g.Header.NumPlayers {
	case 2:
		g.twoPlayerDeal()
	default:
		g.normalDeal()
	}

	g.NewEntry(dealTemplate, sn.H{
		"PID":        g.dealer().id(),
		"HandNumber": g.currentHand(),
	})
}

// assumes deck is shuffled
func (g *game) normalDeal() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	const handSize = 13

	// splits deck into chunks of 13 cards, with possible remainder left for deck
	hands := pie.Chunk(g.State.Deck, handSize)

	// assign chunks/hands to players
	for i := range g.Players {
		g.Players[i].Hand = hands[i]
	}

	// if a remainder chunk, assign to deck
	if len(hands) > g.Header.NumPlayers {
		g.State.Deck = hands[g.Header.NumPlayers]
	}
}

// assumes deck is shuffled
func (g *game) twoPlayerDeal() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Players[0].Hand = g.State.Deck[0:6]
	g.Players[1].Hand = g.State.Deck[6:12]

	g.Players[0].Stack0 = g.State.Deck[12:14]
	g.Players[1].Stack0 = g.State.Deck[14:16]

	g.Players[0].Stack1 = g.State.Deck[16:18]
	g.Players[1].Stack1 = g.State.Deck[18:20]

	g.Players[0].Stack2 = g.State.Deck[20:22]
	g.Players[1].Stack2 = g.State.Deck[22:24]

	g.Players[0].Stack3 = g.State.Deck[24:26]
	g.Players[1].Stack3 = g.State.Deck[26:28]

	g.Players[0].Stack4 = g.State.Deck[28:30]
	g.Players[1].Stack4 = g.State.Deck[30:32]

	g.Players[0].Stack0[0].FaceUp = false
	g.Players[1].Stack0[0].FaceUp = false

	g.Players[0].Stack1[0].FaceUp = false
	g.Players[1].Stack1[0].FaceUp = false

	g.Players[0].Stack2[0].FaceUp = false
	g.Players[1].Stack2[0].FaceUp = false

	g.Players[0].Stack3[0].FaceUp = false
	g.Players[1].Stack3[0].FaceUp = false

	g.Players[0].Stack4[0].FaceUp = false
	g.Players[1].Stack4[0].FaceUp = false

	g.State.Deck = g.State.Deck[32:34]
}
