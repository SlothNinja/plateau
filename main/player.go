package main

import (
	"fmt"
	"sort"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

// Player represents a player of the game.
type player struct {
	ID              sn.PID
	Passed          bool
	Bid             bool
	Colors          []sn.Color
	PerformedAction bool
	Hand            []card
	sn.Stats
}

func (p *player) copy() *player {
	cs := make([]sn.Color, len(p.Colors))
	copy(cs, p.Colors)

	h := make([]card, len(p.Hand))
	copy(h, p.Hand)

	return &player{
		ID:              p.ID,
		Passed:          p.Passed,
		Bid:             p.Bid,
		Colors:          cs,
		PerformedAction: p.PerformedAction,
		Stats:           p.Stats,
		Hand:            h,
	}
}

func (p *player) reset() {
	p.PerformedAction = false
}

func (p *player) bidReset() {
	p.Bid = false
	p.Passed = false
}

func sortedByScore(ps []*player, c comparison) {
	sort.SliceStable(ps, func(i, j int) bool { return ps[i].compareByScore(ps[j]) == c })
}

func (p *player) compareByScore(p2 *player) comparison {
	switch {
	case p.Score < p2.Score:
		return lessThan
	case p.Score > p2.Score:
		return greaterThan
	default:
		return equalTo
	}
}

type comparison int64

const (
	equalTo     comparison = 0
	lessThan    comparison = -1
	greaterThan comparison = 1
	ascending              = lessThan
	descending             = greaterThan
)

// equal returns true if players equal, false otherwise.
func (p *player) equal(op *player) bool {
	return p == op
}

func (client Client) determinePlaces(g *game) (sn.Results, error) {
	rs := make(sn.Results)
	sortedByScore(g.Players, descending)
	ps := make([]*player, len(g.Players))
	cnt := copy(ps, g.Players)
	if cnt != len(g.Players) {
		return nil, fmt.Errorf("error copying players")
	}

	place := 1
	for len(ps) != 0 {
		// Find all players tied at place
		found := pie.Filter(ps, func(p *player) bool { return ps[0].compareByScore(p) == equalTo })
		// Get user keys for found players
		rs[place] = pie.Map(found, func(p *player) *datastore.Key { return g.userKeyFor(p.ID) })
		// Set ps to remaining players
		_, ps = pie.Diff(ps, found)
		// Above does not guaranty order so sort
		sortedByScore(ps, descending)
		// Increase place by number of players added to current place
		place += len(rs[place])
	}
	// 	for _, p1 := range g.players {
	// 		rs[place] = append(rs[place], g.userKeyFor(p1.ID))
	// 		for _, p2 := range g.players {
	// 			if p1.ID != p2.ID && p1.compare(p2) != equalTo {
	// 				place++
	// 				break
	// 			}
	// 		}
	// 	}
	return rs, nil
}

func (g *game) addNewPlayers() {
	g.Players = make([]*player, g.NumPlayers)
	for i := range g.Players {
		g.Players[i] = g.newPlayer(i)
	}
}

func (g game) newPlayer(i int) *player {
	return &player{
		ID:     sn.PID(i + 1),
		Colors: defaultColors(),
	}
}

// IndexFor returns the index for the player and bool indicating whether player found.
// if not found, returns -1
func (g game) indexFor(p1 *player) int {
	return pie.FindFirstUsing(g.Players, func(p2 *player) bool { return p1.equal(p2) })
}

func (g game) playerByPID(pid sn.PID) *player {
	const notFound = -1
	index := pie.FindFirstUsing(g.Players, func(p *player) bool { return p.ID == pid })
	if index == notFound {
		return nil
	}
	return g.Players[index]
}

func (g game) playerByUID(uid sn.UID) *player {
	index := sn.UIndex(pie.FindFirstUsing(g.UserIDS, func(id sn.UID) bool { return id == uid }))
	return g.playerByPID(index.ToPID())
}

// treats players as a circular buffer, thus permitting indices larger than length and indices less than 0
func (g game) playerByIndex(i int) *player {
	l := len(g.Players)
	if l < 1 {
		return nil
	}

	r := i % l
	if r < 0 {
		return g.Players[l+r]
	}
	return g.Players[r]
}

func pidsFor(ps []*player) []sn.PID {
	return pie.Map(ps, func(p *player) sn.PID { return p.ID })
}

func (g game) uidForPID(pid sn.PID) sn.UID {
	return g.UserIDS[pid.ToIndex()]
}

func (g game) playerUIDS() []sn.UID {
	return pie.Map(g.Players, func(p *player) sn.UID { return g.UserIDS[p.ID.ToIndex()] })
}

func (g game) playerStats() []sn.Stats {
	return pie.Map(g.Players, func(p *player) sn.Stats { return p.Stats })
}

func (g game) uidsForPIDS(pids []sn.PID) []sn.UID {
	return pie.Map(pids, func(pid sn.PID) sn.UID { return g.uidForPID(pid) })
}

func (g game) userKeyFor(pid sn.PID) *datastore.Key {
	return sn.NewUser(g.uidForPID(pid)).Key
}

type test func(*player) bool

// cp specifies the current
// return player after cp that satisfies all tests ts
// if tests ts is empty, return player after cp
func (g game) nextPlayer(cp *player, ts ...test) *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	start := g.indexFor(cp) + 1
	stop := start + g.NumPlayers

	// g.playerByIndex uses the players slice as-if it were circular buffer
	// start is one index after cp
	// stop is num players later, thus one pass through circular buffer
	for i := start; i <= stop; i++ {
		np := g.playerByIndex(i)
		if pie.All(ts, func(t test) bool { return t(np) }) {
			return np
		}
	}
	return nil
}

func (g game) declarer() *player {
	return g.playerByPID(pie.First(g.DeclarersTeam))
}

func (g game) partner() *player {
	if len(g.DeclarersTeam) == 2 {
		return g.playerByPID(pie.Last(g.DeclarersTeam))
	}
	return nil
}

func (g game) partners() []*player {
	if len(g.DeclarersTeam) < 2 {
		return nil
	}
	return pie.Map(g.DeclarersTeam[1:], func(pid sn.PID) *player { return g.playerByPID(pid) })
}

func (g game) declarers() []*player {
	return pie.Filter(g.Players, func(p *player) bool {
		return pie.Any(g.DeclarersTeam, func(pid sn.PID) bool {
			return pid == p.ID
		})
	})
}

func (g game) opposersTeam() []sn.PID {
	return pie.FilterNot(pidsFor(g.Players), func(pid1 sn.PID) bool {
		return pie.Any(g.DeclarersTeam, func(pid2 sn.PID) bool {
			return pid1 == pid2
		})
	})
}

func (g game) opposers() []*player {
	return pie.FilterNot(g.Players, func(p *player) bool {
		return pie.Any(g.DeclarersTeam, func(pid sn.PID) bool {
			return pid == p.ID
		})
	})
}
