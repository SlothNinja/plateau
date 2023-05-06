package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"gonum.org/v1/gonum/graph/path"
)

type trick struct {
	Cards []card
	WonBy sn.PID
}

// sn.Header.Turn used to track zero based index of trick slice
func (g game) trickIndex() int {
	return g.Turn
}

func (g *game) nextTrickIndex() int {
	g.Turn++
	return g.Turn
}

func (g *game) resetTrickIndex() int {
	g.Turn = 0
	return g.Turn
}

func (g game) trickNumber() int {
	return g.trickIndex() + 1
}

func (g game) currentTrick() trick {
	return (g.Tricks[g.trickIndex()])
}

func (g *game) endTrick() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	ledSuit := g.ledSuit()
	var winningCard card

	for _, c := range g.currentTrick().Cards {
		if ((c.Suit == ledSuit) || (c.Suit == trumps)) && c.value() > winningCard.value() {
			winningCard = c
		}
	}

	g.Tricks[g.trickIndex()].WonBy = winningCard.PlayedBy
	g.nextTrickIndex()

	return g.playerByPID(winningCard.PlayedBy)
}

func (g game) allCardsPlayed() bool {
	return pie.All(g.Players, func(p *player) bool { return len(p.Hand) == 0 })
}

func (g game) objectiveMade() ([]node, bool) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return g.objectiveTest(g.spacesFor(g.DeclarersTeam))
}

func (g game) objectiveBlocked() ([]node, bool) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	path, found := g.objectiveTest(g.spacesNotOwnedBy(g.opposersTeam()))
	return path, !found
}

func (g game) objectiveTest(ss []space) ([]node, bool) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	graph := g.graphFor(ss)
	paths := path.DijkstraAllPaths(graph)
	switch g.currentBid().Objective {
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
	return pie.Filter(g.Tricks, func(t trick) bool { return pie.Contains(team, t.WonBy) })
}

func (g *game) trickWonBy(team []sn.PID) (won []bool) {
	pie.Each(g.Tricks, func(t trick) {
		won = append(won, pie.Contains(team, t.WonBy))
	})
	return won
}
