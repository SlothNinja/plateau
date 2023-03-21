package main

import (
	"fmt"
	"sort"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v2"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

// Player represents a player of the game.
type player struct {
	ID              sn.PID     `json:"id"`
	Score           int        `json:"score"`
	Passed          bool       `json:"passed"`
	Colors          []sn.Color `json:"colors"`
	PerformedAction bool       `json:"performedAction"`
	Stats           stats      `json:"stats"`
	HasBid          bool       `json:"hasBid"`
	// 	Chips            chips        `json:"chips"`
	// 	PlayedChips      chips        `json:"playedChips"`
	// 	Office           office       `json:"office"`
	// 	SlanderChips     slanderChips `json:"slanderChips"`
	// 	PlacedBosses     int          `json:"placedBosses"`
	// 	PlacedImmigrants int          `json:"placedImmigrants"`
	// 	LockedUp         int          `json:"lockedUp"`
	// 	Slandered        int          `json:"slandered"`
	// 	Candidate        bool         `json:"candidate"`
	// 	UsedOffice       bool         `json:"usedOffice"`
}

func (p *player) reset() {
	p.PerformedAction = false
	p.HasBid = false
	// p.PlacedBosses = 0
	// p.PlacedImmigrants = 0
	// p.Slandered = 0
	// p.Candidate = false
	// p.UsedOffice = false
}

func (g *game) beginningOfTurnReset() {
	// g.slanderedPlayerID = sn.NoPID
	// g.slanderNationality = noNationality
	// g.currentWardID = noWardID
	// g.moveFromWardID = noWardID
}

// type slanderChips map[int]bool

// func sortedByAll(ps []*player, c comparison) {
// 	sort.SliceStable(ps, func(i, j int) bool { return ps[i].compare(ps[j]) == c })
// }

// func sortedByChipsAndMayor(ps []*player, c comparison) {
// 	sort.SliceStable(ps, func(i, j int) bool { return ps[i].compareWithoutScore(ps[j]) == c })
// }

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

// func (p *player) compare(p2 *player) comparison {
// 	c := p.compareByScore(p2)
// 	if c == equalTo {
// 		return p.compareWithoutScore(p2)
// 	}
// 	return c
// }

// func (p *player) compareWithoutScore(p2 *player) comparison {
// 	c := p.compareByTotalChips(p2)
// 	if c != equalTo {
// 		return c
// 	}
// 	c = p.compareByFavors(p2)
// 	if c != equalTo {
// 		return c
// 	}
// 	return p.compareByMayor(p2)
// }

// func (p *player) compareByTotalChips(p2 *player) comparison {
// 	switch {
// 	case p.Chips.count() < p2.Chips.count():
// 		return lessThan
// 	case p.Chips.count() > p2.Chips.count():
// 		return greaterThan
// 	default:
// 		return equalTo
// 	}
// }

// func (p *player) compareByFavors(p2 *player) comparison {
// 	for _, n := range nationalities() {
// 		switch {
// 		case p.Chips[n] < p2.Chips[n]:
// 			return lessThan
// 		case p.Chips[n] > p2.Chips[n]:
// 			return greaterThan
// 		}
// 	}
// 	return equalTo
// }
//
// func (p *player) compareByMayor(p2 *player) comparison {
// 	m1, m2 := p.Office == mayor, p2.Office == mayor
// 	switch {
// 	case !m1 && m2:
// 		return lessThan
// 	case m1 && !m2:
// 		return greaterThan
// 	default:
// 		return equalTo
// 	}
// }

// equal returns true if players equal, false otherwise.
func (p *player) equal(op *player) bool {
	return p != nil && op != nil && p.ID == op.ID
}

func (client *Client) determinePlaces(c *gin.Context, g *game) (sn.Results, error) {
	rs := make(sn.Results)
	sortedByScore(g.players, descending)
	ps := make([]*player, len(g.players))
	cnt := copy(ps, g.players)
	if cnt != len(g.players) {
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

// // chips stores favor chips by nationality.
// type chips map[nationality]int
//
// // Count provides a total count of favor chips regardless of nationality.
// func (cs chips) count() int {
// 	var cnt int
// 	for _, value := range cs {
// 		cnt += value
// 	}
// 	return cnt
// }

// // maxInfluenceIn returns the maximal amount of influence a player has in a ward w if bid all relevant favor chips.
// func (p *player) maxInfluenceIn(w *ward) int {
// 	return w.bossesFor(p) + w.playableChipsFor(p)
// }
//
// func (p *player) electionCountIn(w *ward) int {
// 	return p.PlayedChips.count() + w.bossesFor(p)
// }
//
// // chipsFor returns the number of favor chips a player has for a given nationality.
// func (p *player) chipsFor(n nationality) int {
// 	return p.Chips[n]
// }

func (g *game) addNewPlayers() {
	g.players = make([]*player, g.NumPlayers)
	for i := range g.players {
		g.players[i] = g.newPlayer(i)
	}
}

func (g *game) newPlayer(i int) *player {
	return &player{
		ID: sn.PID(i + 1),
		// Colors: defaultColors()[:g.NumPlayers],
		// Chips:        make(chips),
		// PlayedChips:  make(chips),
		// SlanderChips: slanderChips{2: true, 3: true, 4: true},
	}
}

// Controlled provides a count of the immigrants of the given nationality in the wards having a boss of the player.
// func (g *game) controlledBy(p *player, n nationality) int {
// 	var cnt int
// 	for _, w := range g.activeWards() {
// 		if w.bossesFor(p) > 0 {
// 			cnt += w.Immigrants[n]
// 		}
// 	}
// 	return cnt
// }
//
// func (p *player) placedPieces() int {
// 	return p.PlacedBosses + p.PlacedImmigrants
// }

// IndexFor returns the index for the player and bool indicating whether player found.
func (g *game) indexFor(p1 *player) (int, bool) {
	const notFound = -1
	const found = true
	index := pie.FindFirstUsing(g.players, func(p2 *player) bool { return p1.equal(p2) })
	if index == notFound {
		return index, !found
	}
	return index, found
}

func (g *game) playerByPID(pid sn.PID) *player {
	const notFound = -1
	index := pie.FindFirstUsing(g.players, func(p *player) bool { return p.ID == pid })
	if index == notFound {
		return nil
	}
	return g.players[index]
}

// treats players as a circular buffer, thus permitting indices larger than length and indices less than 0
func (g *game) playerByIndex(i int) *player {
	l := len(g.players)
	if l < 1 {
		return nil
	}

	r := i % l
	if r < 0 {
		return g.players[l+r]
	}
	return g.players[r]
}

func (g *game) playerByUserKey(key *datastore.Key) *player {
	const notFound = -1
	i := pie.FindFirstUsing(g.UserKeys, func(k *datastore.Key) bool { return key.Equal(k) })
	if i == notFound {
		return nil
	}
	return g.players[i]
}

func pids(ps []*player) []sn.PID {
	return pie.Map(ps, func(p *player) sn.PID { return p.ID })
}

func (g *game) uidForPID(pid sn.PID) int64 {
	return g.UserIDS[pid.ToIndex()]
}

func (g *game) uidsForPIDS(pids []sn.PID) []int64 {
	return pie.Map(pids, func(pid sn.PID) int64 { return g.uidForPID(pid) })
}

func (g *game) userKeyFor(pid sn.PID) *datastore.Key {
	return g.UserKeys[pid.ToIndex()]
}
