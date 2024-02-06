package client

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/barkimedes/go-deepcopy"
	"github.com/elliotchance/pie/v2"
)

const endHandTemplate = "end-hand-results"

func (g *game) startEndHandPhase(result handResult, path []node) *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Header.Phase = endHandPhase
	scored := g.scoreHand(result)
	results := g.saveLastResult(result, path, scored)

	g.NewEntry(endHandTemplate, sn.H{"Results": results})

	if end := g.endGameCheck(); end {
		return nil
	}
	return g.startHand()
}

func (g *game) endGameCheck() bool {
	return g.currentHand() == g.finalHand()
}

func (g *game) endHandCheck() (end bool, result handResult, path []node) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	if g.allCardsPlayed() {
		g.revealTalon()
		_, result, path = g.objectiveCheck()
		return true, result, path
	}
	return g.objectiveCheck()
}

func (g *game) objectiveCheck() (end bool, result handResult, path []node) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	path, result = g.objectiveMade()
	if result == dSuccess {
		return true, result, path
	}

	_, result = g.objectiveBlocked()
	if result == dFail {
		return true, result, path
	}
	return false, dPush, nil
}

type handResult string

const (
	dPush    handResult = "push"
	dSuccess handResult = "success"
	dFail    handResult = "failure"
)

func (g *game) scoreHand(result handResult) []int64 {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	if result == dPush {
		return make([]int64, len(g.Players))
	}

	oldscores := pie.Map(g.Players, func(p *player) int64 { return p.Score })
	dtbv := g.currentBidValue()
	dl := len(g.State.DeclarersTeam)
	ol := g.Header.NumPlayers - dl
	if result == dFail {
		dtbv = -dtbv
	}

	pie.Each(g.opposers(), func(p *player) { p.Score -= dtbv })
	switch {
	case dl == ol:
		pie.Each(g.declarers(), func(p *player) { p.Score += dtbv })
	case dl == 1:
		g.declarer().Score += dtbv * int64(ol)
	case dl == 2 && g.Header.NumPlayers == 5:
		g.declarer().Score += dtbv * 2
		g.partner().Score += dtbv
	case dl == 2 && g.Header.NumPlayers == 6:
		pie.Each(g.declarers(), func(p *player) { p.Score += dtbv * 2 })
	}
	deltaScores := make([]int64, len(g.Players))
	for i, p := range g.Players {
		deltaScores[i] = p.Score - oldscores[i]
	}
	return deltaScores
}

func (g *game) startHand() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.State.Deck = deckFor(g.Header.NumPlayers)
	g.Header.Phase = dealPhase

	g.resetTrickIndex()
	g.nextHand()

	if g.currentHand() == 1 {
		g.RandomizePlayers()
	} else {
		g.newDealer()
	}

	switch g.Header.NumPlayers {
	case 2:
		g.State.Tricks = make([]trick, 17) // 16 trick plus talon
	case 6:
		g.State.Tricks = make([]trick, 13) // 13 trick no talon
	default:
		g.State.Tricks = make([]trick, 14) // 13 trick plus talon
	}

	g.State.DeclarersTeam = nil

	g.deal()
	return g.startBidPhase()
}

func (g *game) revealTalon() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	var talon trick
	talon.Cards = g.State.Deck
	if g.currentBid().Exchange == noExchangeBid {
		talon.WonBy = pie.First(g.State.DeclarersTeam)
	} else {
		talon.WonBy = pie.First(g.opposersTeam())
	}
	talonIndex := len(g.State.Tricks) - 1
	g.State.Tricks[talonIndex] = talon
}

type lastResult struct {
	HandNumber    int
	Bids          []bid
	SeatOrder     []sn.PID
	DeclarersTeam []sn.PID
	Tricks        []trick
	Path          []node
	Scored        []int64
	Success       handResult
}

func (g *game) saveLastResult(result handResult, path []node, scored []int64) lastResult {
	last := lastResult{
		HandNumber:    g.currentHand(),
		Bids:          g.State.Bids,
		SeatOrder:     g.Header.OrderIDS,
		DeclarersTeam: g.State.DeclarersTeam,
		Tricks:        g.State.Tricks,
		Path:          path,
		Scored:        scored,
		Success:       result,
	}
	g.State.LastResults = append(g.State.LastResults, deepcopy.MustAnything(last).(lastResult))
	return last
}

func (g *game) currentHand() int {
	return g.Header.Round
}

func (g *game) nextHand() int {
	g.Header.Round++
	return g.Header.Round
}
