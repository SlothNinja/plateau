package main

import (
	"encoding/json"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"gonum.org/v1/gonum/graph/path"
)

type trick struct {
	cards []card
	wonBy sn.PID
}

type jTrick struct {
	Cards []card `json:"cards"`
	WonBy sn.PID `json:"wonBy"`
}

func (t trick) MarshalJSON() ([]byte, error) {
	return json.Marshal(jTrick{
		Cards: t.cards,
		WonBy: t.wonBy,
	})
}

func (t *trick) UnmarshalJSON(bs []byte) error {
	obj := new(jTrick)
	err := json.Unmarshal(bs, obj)
	if err != nil {
		return err
	}

	t.cards = obj.Cards
	t.wonBy = obj.WonBy
	return nil
}

func (g *game) endTrick() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	ledSuit := g.ledSuit()
	var winningCard card

	// warning card variable shadows card type in for loop
	for _, card := range g.tricks[g.trickNumber()].cards {
		if ((card.suit == ledSuit) || (card.suit == trumps)) && card.value() > winningCard.value() {
			winningCard = card
		}
	}

	g.tricks[g.trickNumber()].wonBy = winningCard.playedBy
	g.incTrickNumber()

	switch {
	case g.madeObjective():
		return g.endHand()
	case g.objectiveBlocked():
		return g.endHand()
	default:
		return g.playerByPID(winningCard.playedBy)
	}
}

func (g game) allCardsPlayed() bool {
	return pie.All(g.players, func(p *player) bool { return len(p.hand) == 0 })
}

func (g game) madeObjective() bool {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return g.objectiveTest(own(g, g.declarersTeam))
}

func (g game) objectiveBlocked() bool {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return !g.objectiveTest(mayOwn(g, g.declarersTeam))
}

func (g game) objectiveTest(t spaceTest) bool {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	graph := g.graphFor(t)
	paths := path.DijkstraAllPaths(graph)
	s1, s2, s3, s4, s5, s6 := side1(graph), side2(graph), side3(graph),
		side4(graph), side5(graph), side6(graph)
	switch g.currentBid().objective {
	case bridgeBid:
		return bridge(paths, s1, s2, s3, s4, s5, s6)
	case yBid:
		return y(paths, s1, s2, s3, s4, s5, s6)
	case forkBid:
		return fork(paths, s1, s2, s3, s4, s5, s6)
	}
	return false
}
