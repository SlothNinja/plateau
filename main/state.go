package main

import (
	"github.com/SlothNinja/sn/v3"
)

// state stores the game state of a Tammany Hall game.
type state struct {
	Players       []*player
	Deck          []card
	DeclarersTeam []sn.PID
	Tricks        []trick
	Bids          []bid
	LastResults   []lastResult
	Pick          []card
}

func (s state) copy() state {
	ps := make([]*player, len(s.Players))
	for i, p := range s.Players {
		ps[i] = p.copy()
	}

	d := make([]card, len(s.Deck))
	copy(d, s.Deck)

	dt := make([]sn.PID, len(s.DeclarersTeam))
	copy(dt, s.DeclarersTeam)

	ts := make([]trick, len(s.Tricks))
	copy(ts, s.Tricks)

	bs := make([]bid, len(s.Bids))
	copy(bs, s.Bids)

	lrs := make([]lastResult, len(s.LastResults))
	for i, l := range s.LastResults {
		lrs[i] = l.copy()
	}

	pc := make([]card, len(s.Pick))
	copy(pc, s.Pick)

	return state{
		Players:       ps,
		Deck:          d,
		DeclarersTeam: dt,
		Tricks:        ts,
		Bids:          bs,
		LastResults:   lrs,
		Pick:          pc,
	}
}
