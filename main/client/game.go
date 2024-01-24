package client

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/barkimedes/go-deepcopy"
	"github.com/elliotchance/pie/v2"
)

type game struct {
	sn.Game[state, player, *player]
}

func newGame() *game {
	return &game{sn.Game[state, player, *player]{}}
}

func (g *game) Start(h sn.Header) sn.PID {
	g.Game.Start(h)
	sn.Debugf("g.Log: %#v", g.Log)
	cp := g.startHand()
	g.SetCurrentPlayers(cp.id())
	return cp.id()
}

func (g *game) dealer() *player {
	return pie.First(g.Players)
}

func (g *game) forehand() *player {
	return pie.First(pie.DropTop(g.Players, 1))
}

// Circular shift left of players so dealer is always first element in slice
func (g *game) newDealer() {
	g.Players = pie.Rotate(g.Players, -1)
	g.UpdateOrder()
}

func (g *game) Views() ([]sn.UID, []*game) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	uids, games := make([]sn.UID, g.Header.NumPlayers+1), make([]*game, g.Header.NumPlayers+1)
	for i, p := range g.Players {
		sn.Debugf("p: %#v", p)
		uids[i] = g.UIDForPID(p.id())
		games[i] = g.viewFor(p)
	}

	// add view for non-player
	uids[g.Header.NumPlayers] = 0
	games[g.Header.NumPlayers] = g.viewFor(nil)

	return uids, games
}

// remove hand of other players and deck from data viewed by player
func (g *game) viewFor(p *player) *game {
	g2 := deepcopy.MustAnything(g).(*game)
	for _, p2 := range g2.Players {
		if p == nil || (p.id() != p2.id()) {
			p2.Hand = nil
		}
		stacksView(p2)
	}
	g2.State.Deck = nil
	return g2
}

func stacksView(p *player) {
	stackView(p.Stack0)
	stackView(p.Stack1)
	stackView(p.Stack2)
	stackView(p.Stack3)
	stackView(p.Stack4)
}

func stackView(stack []card) {
	for i, c := range stack {
		if !c.FaceUp {
			stack[i].Rank = noRank
			stack[i].Suit = noSuit
		}
	}
}
