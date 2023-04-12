package main

import (
	"encoding/json"
	"fmt"
	"sort"

	"cloud.google.com/go/datastore"
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

// Player represents a player of the game.
type player struct {
	id              sn.PID     // `json:"id"`
	score           int        // `json:"score"`
	passed          bool       // `json:"passed"`
	bid             bool       // `json:"bid"`
	colors          []sn.Color // `json:"colors"`
	performedAction bool       // `json:"performedAction"`
	stats           stats      // `json:"stats"`
	hand            []card     // `json:"hand"`
}

type jPlayer struct {
	ID              sn.PID     `json:"id"`
	Score           int        `json:"score"`
	Passed          bool       `json:"passed"`
	Bid             bool       `json:"bid"`
	Colors          []sn.Color `json:"colors"`
	PerformedAction bool       `json:"performedAction"`
	Stats           stats      `json:"stats"`
	Hand            []card     `json:"hand"`
}

func (p player) MarshalJSON() ([]byte, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	return json.Marshal(jPlayer{
		ID:              p.id,
		Score:           p.score,
		Passed:          p.passed,
		Bid:             p.bid,
		Colors:          p.colors,
		PerformedAction: p.performedAction,
		Stats:           p.stats,
		Hand:            p.hand,
	})
}

func (p *player) UnmarshalJSON(bs []byte) error {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var obj jPlayer
	err := json.Unmarshal(bs, &obj)
	if err != nil {
		return err
	}
	p.id = obj.ID
	p.score = obj.Score
	p.passed = obj.Passed
	p.bid = obj.Bid
	p.colors = obj.Colors
	p.performedAction = obj.PerformedAction
	p.stats = obj.Stats
	p.hand = obj.Hand
	return nil
}

func (p *player) reset() {
	p.performedAction = false
	// p.PlacedBosses = 0
	// p.PlacedImmigrants = 0
	// p.Slandered = 0
	// p.Candidate = false
	// p.UsedOffice = false
}

func (p *player) bidReset() {
	p.bid = false
	p.passed = false
}

// func (g *game) beginningOfTurnReset() {
// 	// g.slanderedPlayerID = sn.NoPID
// 	// g.slanderNationality = noNationality
// 	// g.currentWardID = noWardID
// 	// g.moveFromWardID = noWardID
// }

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
	case p.score < p2.score:
		return lessThan
	case p.score > p2.score:
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
	return p != nil && op != nil && p.id == op.id
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
		rs[place] = pie.Map(found, func(p *player) *datastore.Key { return g.userKeyFor(p.id) })
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
		id:     sn.PID(i + 1),
		colors: defaultColors(),
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
// if not found, returns -1
func (g *game) indexFor(p1 *player) int {
	return pie.FindFirstUsing(g.players, func(p2 *player) bool { return p1.equal(p2) })
}

func (g *game) playerByPID(pid sn.PID) *player {
	const notFound = -1
	index := pie.FindFirstUsing(g.players, func(p *player) bool { return p.id == pid })
	if index == notFound {
		return nil
	}
	return g.players[index]
}

func (g *game) playerByUID(uid sn.UID) *player {
	const notFound = -1
	index := pie.FindFirstUsing(g.UserIDS, func(id sn.UID) bool { return id == uid })
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

func pidsFor(ps []*player) []sn.PID {
	return pie.Map(ps, func(p *player) sn.PID { return p.id })
}

func (g *game) uidForPID(pid sn.PID) sn.UID {
	return g.UserIDS[pid.ToIndex()]
}

func (g *game) uidsForPIDS(pids []sn.PID) []sn.UID {
	return pie.Map(pids, func(pid sn.PID) sn.UID { return g.uidForPID(pid) })
}

func (g *game) userKeyFor(pid sn.PID) *datastore.Key {
	return g.UserKeys[pid.ToIndex()]
}

// // ps is an optional parameter.
// // If no player is provided, assume current player.
// func (g *game) nextPlayer(ps ...*player) *player {
// 	if len(ps) == 1 {
// 		i, _ := g.indexFor(ps[0])
// 		return g.playerByIndex(i + 1)
// 	}
//
// 	cp := g.currentPlayer()
// 	if cp == nil {
// 		return nil
// 	}
//
// 	i, _ := g.indexFor(cp)
// 	return g.playerByIndex(i + 1)
// }

type test func(*player) bool

// cp specifies the current
// return player after cp that satisfies all tests ts
// if tests ts is empty, return player after cp
func (g game) nextPlayer(cp *player, ts ...test) (np *player) {
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
	return g.playerByPID(pie.First(g.declarersTeam))
}

func (g game) partner() *player {
	if len(g.declarersTeam) == 2 {
		return g.playerByPID(pie.Last(g.declarersTeam))
	}
	return nil
}

func (g game) partners() []*player {
	if len(g.declarersTeam) < 2 {
		return nil
	}
	return pie.Map(g.declarersTeam[1:], func(pid sn.PID) *player { return g.playerByPID(pid) })
}

func (g game) declarers() []*player {
	return pie.Filter(g.players, func(p *player) bool {
		return pie.Any(g.declarersTeam, func(pid sn.PID) bool {
			return pid == p.id
		})
	})
}

func (g game) opposersTeam() []sn.PID {
	return pie.FilterNot(pidsFor(g.players), func(pid1 sn.PID) bool {
		return pie.Any(g.declarersTeam, func(pid2 sn.PID) bool {
			return pid1 == pid2
		})
	})
}

func (g game) opposers() []*player {
	return pie.FilterNot(g.players, func(p *player) bool {
		return pie.Any(g.declarersTeam, func(pid sn.PID) bool {
			return pid == p.id
		})
	})
}
