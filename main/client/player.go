package client

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

// Player represents a player of the game.
type player struct {
	Bid  bool
	Hand []card

	// firestore does not support a slice of slices
	// thus the stacks are not implemented as a slice of stacks
	Stack0 []card
	Stack1 []card
	Stack2 []card
	Stack3 []card
	Stack4 []card

	sn.Player
}

func (p player) stacks() [][]card {
	return [][]card{p.Stack0, p.Stack1, p.Stack2, p.Stack3, p.Stack4}
}

func (p player) playableStacks() []card {
	var cards []card
	for _, stack := range p.stacks() {
		if len(stack) > 0 {
			cards = append(cards, pie.Last(stack))
		}
	}
	return cards
}

func (p player) playableCards() []card {
	return append(p.Hand, p.playableStacks()...)
}

func (p player) hasCard(c card) bool {
	return pie.Contains(p.playableCards(), c)
}

func (p player) hasSuit(s suit) bool {
	return pie.Any(p.playableCards(), func(c card) bool { return c.Suit == s })
}

func (p player) hasRank(r rank) bool {
	return pie.Any(p.playableCards(), func(c card) bool { return c.Rank == r })
}

func (p *player) updateStacks() {
	updateStack(p.Stack0)
	updateStack(p.Stack1)
	updateStack(p.Stack2)
	updateStack(p.Stack3)
	updateStack(p.Stack4)
}

func updateStack(stack []card) {
	if len(stack) == 1 {
		stack[0].FaceUp = true
	}
}

func (p *player) id() sn.PID {
	if p == nil {
		return sn.NoPID
	}
	return p.Player.ID
}

func (p *player) reset() {
	p.PerformedAction = false
}

func (p *player) bidReset() {
	p.Bid = false
	p.Passed = false
}

func (g *game) declarer() *player {
	return g.PlayerByPID(pie.First(g.State.DeclarersTeam))
}

func (g *game) partner() *player {
	if len(g.State.DeclarersTeam) == 2 {
		return g.PlayerByPID(pie.Last(g.State.DeclarersTeam))
	}
	return nil
}

func (g *game) partners() []*player {
	if len(g.State.DeclarersTeam) < 2 {
		return nil
	}
	return pie.Map(g.State.DeclarersTeam[1:], func(pid sn.PID) *player { return g.PlayerByPID(pid) })
}

func (g *game) declarers() []*player {
	return pie.Filter(g.Players, func(p *player) bool {
		return pie.Any(g.State.DeclarersTeam, func(pid sn.PID) bool {
			return pid == p.id()
		})
	})
}

func (g *game) opposersTeam() []sn.PID {
	return pie.FilterNot(g.Players.PIDS(), func(pid1 sn.PID) bool {
		return pie.Any(g.State.DeclarersTeam, func(pid2 sn.PID) bool {
			return pid1 == pid2
		})
	})
}

func (g *game) opposers() []*player {
	return pie.FilterNot(g.Players, func(p *player) bool {
		return pie.Any(g.State.DeclarersTeam, func(pid sn.PID) bool {
			return pid == p.id()
		})
	})
}
