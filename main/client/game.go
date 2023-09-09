package client

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/barkimedes/go-deepcopy"
	"github.com/elliotchance/pie/v2"
)

// Game provides a Le Plateau game.
type game struct {
	// Log glog
	// sn.Header
	sn.Game[*player]
	state
}

func (g *game) Start(h *sn.Header) sn.Playerer {
	g.Game.Start(h)
	sn.Debugf("g.Log: %#v", g.Log)
	cp := g.startHand()
	g.SetCurrentPlayers(cp)
	return cp
}

func (g game) dealer() *player {
	return pie.First(g.Players)
}

func (g game) forehand() *player {
	return pie.First(pie.DropTop(g.Players, 1))
}

// Circular shift left of players so dealer is always first element in slice
func (g *game) newDealer() {
	g.Players = pie.Rotate(g.Players, -1)
	g.UpdateOrder()
}

func (g *game) Views() ([]sn.UID, []*game) {
	uids, games := make([]sn.UID, g.NumPlayers+1), make([]*game, g.NumPlayers+1)
	for i, p := range g.Players {
		uids[i] = g.UIDForPID(p.ID)
		games[i] = g.viewFor(p)
	}

	// add view for non-player
	uids[g.NumPlayers] = 0
	games[g.NumPlayers] = g.viewFor(nil)

	return uids, games
}

// remove hand of other players and deck from data viewed by player
func (g *game) viewFor(p *player) *game {
	g2 := deepcopy.MustAnything(g).(*game)
	for _, p2 := range g2.Players {
		if p == nil || p.ID != p2.ID {
			p2.Hand = nil
		}
		stacksView(p2)
	}
	g2.Deck = nil
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

// // not truly a deep copy, though state is deeply copied.
// func (g game) copy() game {
// 	var g2 game
// 	g2.Header = g.Header
// 	g2.Log = g.Log
// 	g2.state = g.state.copy()
// 	g2.Players = copyPlayers(g.Players)
// 	return g2
// }
//
// func copyPlayers(ps []*player) []*player {
// 	return pie.Map(ps, func(p *player) *player { return p.Copy().(*player) })
// }
