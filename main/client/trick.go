package client

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
func (g *game) trickIndex() int {
	return g.Header.Turn
}

func (g *game) nextTrickIndex() int {
	g.Header.Turn++
	return g.Header.Turn
}

func (g *game) resetTrickIndex() int {
	g.Header.Turn = 0
	return g.Header.Turn
}

func (g *game) trickNumber() int {
	return g.trickIndex() + 1
}

func (g *game) currentTrick() trick {
	return (g.State.Tricks[g.trickIndex()])
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

	for _, p := range g.Players {
		p.updateStacks()
	}
	g.State.Tricks[g.trickIndex()].WonBy = winningCard.PlayedBy
	g.nextTrickIndex()

	return g.PlayerByPID(winningCard.PlayedBy)
}

func (g *game) allCardsPlayed() bool {
	return pie.All(g.Players, func(p *player) bool { return len(p.playableCards()) == 0 })
}

func (g *game) allPassed() bool {
	return pie.All(g.Players, func(p *player) bool { return p.Passed })
}

func (g *game) objectiveMade() ([]node, handResult) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return g.objectiveTest(g.spacesFor(g.State.DeclarersTeam))
}

func (g *game) objectiveBlocked() ([]node, handResult) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	switch path, result := g.objectiveTest(g.spacesNotOwnedBy(g.opposersTeam())); result {
	case dSuccess:
		return path, dSuccess
	case dFail:
		return path, dFail
	default:
		return path, dPush
	}
}

func (g *game) objectiveTest(ss []space) ([]node, handResult) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	graph := g.graphFor(ss)
	paths := path.DijkstraAllPaths(graph)
	switch g.currentBid().Objective {
	case bridgeBid:
		return g.bridge(graph, paths)
	case yBid:
		return g.y(graph, paths)
	case forkBid:
		return g.fork(graph, paths)
	case fiveSidesBid:
		return g.fiveSides(graph, paths)
	case sixSidesBid:
		return g.sixSides(graph, paths)
	default:
		return nil, dFail
	}
}

func (g *game) tricksFor(team []sn.PID) []trick {
	return pie.Filter(g.State.Tricks, func(t trick) bool { return pie.Contains(team, t.WonBy) })
}

func (g *game) trickWonBy(team []sn.PID) (won []bool) {
	pie.Each(g.State.Tricks, func(t trick) {
		won = append(won, pie.Contains(team, t.WonBy))
	})
	return won
}
