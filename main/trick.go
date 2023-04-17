package main

import (
	"encoding/json"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"gonum.org/v1/gonum/graph"
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

// sn.Header.Turn used to track card played in current hand
// zero based to align with indices of trick slice
func (g game) trickNumber() int {
	return g.Turn
}

func (g *game) incTrickNumber() int {
	g.Turn++
	return g.Turn
}

func (g *game) resetTrickNumber() int {
	g.Turn = 0
	return g.Turn
}

func (g game) currentTrick() trick {
	return (g.tricks[g.trickNumber()])
}

func (g *game) endTrick() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	ledSuit := g.ledSuit()
	var winningCard card

	// warning card variable shadows card type in for loop
	for _, card := range g.currentTrick().cards {
		if ((card.suit == ledSuit) || (card.suit == trumps)) && card.value() > winningCard.value() {
			winningCard = card
		}
	}

	g.tricks[g.trickNumber()].wonBy = winningCard.playedBy
	g.incTrickNumber()

	// if _, successful := g.madeObjective(); successful {
	// 	g.startEndHandPhase()
	// 	g.startHand()
	// 	return g.startBidPhase()
	// }

	// if _, blocked := g.objectiveBlocked(); blocked {
	// 	g.startEndHandPhase()
	// 	g.startHand()
	// 	return g.startBidPhase()
	// }

	return g.playerByPID(winningCard.playedBy)
}

func (g game) allCardsPlayed() bool {
	return pie.All(g.players, func(p *player) bool { return len(p.hand) == 0 })
}

func (g game) madeObjective() ([]graph.Node, bool) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return g.objectiveTest(g.spacesFor(g.declarersTeam))
}

func (g game) objectiveBlocked() ([]graph.Node, bool) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return g.objectiveTest(g.spacesNotOwnedBy(g.opposersTeam()))
}

func (g game) objectiveTest(ss []space) ([]graph.Node, bool) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	graph := g.graphFor(ss)
	paths := path.DijkstraAllPaths(graph)
	switch g.currentBid().objective {
	case bridgeBid:
		return bridge(graph, paths)
	case yBid:
		return y(graph, paths)
	case forkBid:
		return fork(graph, paths)
	case fiveSidesBid:
		return fiveSides(graph, paths)
	case sixSidesBid:
		return sixSides(graph, paths)
	default:
		return nil, false
	}
}

func (g *game) tricksFor(team []sn.PID) []trick {
	return pie.Filter(g.tricks, func(t trick) bool { return pie.Contains(team, t.wonBy) })
}

func (g *game) trickWonBy(team []sn.PID) (won []bool) {
	pie.Each(g.tricks, func(t trick) {
		won = append(won, pie.Contains(team, t.wonBy))
	})
	return won
}
