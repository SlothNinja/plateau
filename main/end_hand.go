package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

func (g *game) startEndHandPhase(success bool, path []node) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = endHandPhase
	scored := g.scoreHand(success)
	g.saveLastResult(success, path, scored)
}

func (g *game) endHandCheck() (end bool, success bool, path []node) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	if g.allCardsPlayed() {
		g.revealTalon()
		_, success, path = g.objectiveCheck()
		return true, success, path
	}
	return g.objectiveCheck()
}

func (g game) objectiveCheck() (end bool, success bool, path []node) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	path, success = g.objectiveMade()
	if success {
		return true, true, path
	}

	_, blocked := g.objectiveBlocked()
	if blocked {
		return true, false, path
	}
	return false, false, nil
}

func (g *game) scoreHand(successful bool) []int64 {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	oldscores := pie.Map(g.Players, func(p *player) int64 { return p.Score })
	dtbv := g.currentBidValue()
	dl := len(g.DeclarersTeam)
	ol := g.NumPlayers - dl
	if !successful {
		dtbv = -dtbv
	}

	pie.Each(g.opposers(), func(p *player) { p.Score -= dtbv })
	switch {
	case dl == ol:
		pie.Each(g.declarers(), func(p *player) { p.Score += dtbv })
	case dl == 1:
		g.declarer().Score += dtbv * int64(ol)
	case dl == 2 && g.NumPlayers == 5:
		g.declarer().Score += dtbv * 2
		g.partner().Score += dtbv
	case dl == 2 && g.NumPlayers == 6:
		pie.Each(g.declarers(), func(p *player) { p.Score += dtbv * 2 })
	}
	deltaScores := make([]int64, len(g.Players))
	for i, p := range g.Players {
		deltaScores[i] = p.Score - oldscores[i]
	}
	return deltaScores
}

func (g *game) startHand() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Deck = deckFor(g.NumPlayers)
	g.Phase = dealPhase

	g.resetTrickIndex()
	g.nextHand()

	if g.currentHand() == 1 {
		g.randomSeats()
	} else {
		g.newDealer()
	}

	switch g.NumPlayers {
	case 2:
		g.Tricks = make([]trick, 17) // 16 trick plus talon
	case 6:
		g.Tricks = make([]trick, 13) // 13 trick no talon
	default:
		g.Tricks = make([]trick, 14) // 13 trick plus talon
	}

	g.DeclarersTeam = nil

	g.deal()
}

func (g *game) revealTalon() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var talon trick
	talon.Cards = g.Deck
	if g.currentBid().Exchange == noExchangeBid {
		talon.WonBy = pie.First(g.DeclarersTeam)
	} else {
		talon.WonBy = pie.First(g.opposersTeam())
	}
	g.Tricks = append(g.Tricks, talon)
}

type lastResult struct {
	Bids          []bid
	SeatOrder     []sn.PID
	DeclarersTeam []sn.PID
	Tricks        []trick
	Path          []node
	Scored        []int64
	Success       bool
}

func (g *game) saveLastResult(success bool, path []node, scored []int64) {
	last := lastResult{
		Bids:          g.Bids,
		SeatOrder:     g.OrderIDS,
		DeclarersTeam: g.DeclarersTeam,
		Tricks:        g.Tricks,
		Path:          path,
		Scored:        scored,
		Success:       success,
	}
	g.LastResults = append(g.LastResults, last.copy())
}

func (l lastResult) copy() lastResult {
	last := lastResult{
		Bids:          make([]bid, len(l.Bids)),
		SeatOrder:     make([]sn.PID, len(l.SeatOrder)),
		DeclarersTeam: make([]sn.PID, len(l.DeclarersTeam)),
		Tricks:        make([]trick, len(l.Tricks)),
		Path:          make([]node, len(l.Path)),
		Scored:        make([]int64, len(l.Scored)),
		Success:       l.Success,
	}
	copy(last.Bids, l.Bids)
	copy(last.SeatOrder, l.SeatOrder)
	copy(last.DeclarersTeam, l.DeclarersTeam)
	copy(last.Tricks, l.Tricks)
	copy(last.Path, l.Path)
	copy(last.Scored, l.Scored)
	return last
}
