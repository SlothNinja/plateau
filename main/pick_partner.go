package main

import (
	"fmt"

	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
)

func (g *game) startPickPartner() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = pickPartnerPhase
	switch b := g.lastBid(); b.Teams {
	case duoBid:
		g.updatePickCards()
		return g.declarer()
	case trioBid:
		g.updatePickCards()
		return g.declarer()
	default:
		return nil
	}
}

func (g *game) updatePickCards() {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	for _, r := range []rank{roiRank, dameRank, cavalierRank, valetRank, tenRank} {
		g.Pick = removeCards(cardsOfRank(r), g.Deck...)
		for _, p := range g.declarers() {
			g.Pick = removeCards(g.Pick, p.Hand...)
		}
		if len(g.Pick) > 0 {
			break
		}
	}
}

func (g game) otherTeam(pids1 []sn.PID) []sn.PID {
	return pie.FilterNot(g.Players.PIDS(), func(pid2 sn.PID) bool {
		return pie.Any(pids1, func(pid1 sn.PID) bool { return pid1 == pid2 })
	})
}

// func (cl Client) pickPartnerHandler(ctx *gin.Context) {
// 	cl.Log.Debugf(msgEnter)
// 	defer cl.Log.Debugf(msgExit)
//
// 	cu, err := cl.Current(ctx)
// 	if err != nil {
// 		cl.Log.Warningf(err.Error())
// 	}
//
// 	g, err := cl.GetGame(ctx, cu)
// 	if err != nil {
// 		sn.JErr(ctx, err)
// 		return
// 	}
//
// 	cp, np, err := g.pickPartner(ctx, cu)
// 	if err != nil {
// 		sn.JErr(ctx, err)
// 		return
// 	}
//
// 	cp.Stats.Moves++
// 	cp.Stats.Think += time.Since(g.UpdatedAt)
//
// 	np.reset()
// 	g.SetCurrentPlayers(np)
//
// 	if err := cl.Commit(ctx, g, cu); err != nil {
// 		sn.JErr(ctx, err)
// 		return
// 	}
//
// 	ctx.JSON(http.StatusOK, nil)
// }

func (g *game) pickPartner(ctx *gin.Context, cu sn.User) (*player, *player, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	// warning card variable shadows card type
	cp, card, err := g.validatePickPartner(ctx, cu)
	if err != nil {
		return cp, nil, err
	}

	for _, p := range g.opposers() {
		if p.hasCard(card) {
			g.DeclarersTeam = append(g.DeclarersTeam, p.ID)
			break
		}
	}

	if g.lastBid().Teams == trioBid && len(g.DeclarersTeam) < 3 {
		g.updatePickCards()
		return cp, cp, nil
	}
	g.Pick = nil
	return cp, g.startIncObjective(), nil
}

func (g game) validatePickPartner(ctx *gin.Context, cu sn.User) (*player, card, error) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	noCard := card{}
	cp, err := g.ValidatePlayerAction(cu)
	if err != nil {
		return nil, noCard, err
	}

	cards, err := getCards(ctx)
	if err != nil {
		return nil, noCard, err
	}

	if len(cards) != 1 {
		return nil, noCard, fmt.Errorf("must select one and only one card: %w", sn.ErrValidation)
	}

	selectedCard := pie.First(cards)

	switch {
	case g.Phase != pickPartnerPhase:
		return nil, noCard, fmt.Errorf("cannot select partner in %q phase: %w", g.Phase, sn.ErrValidation)
	case cp.hasCard(selectedCard):
		return nil, noCard, fmt.Errorf("must select highest non-trump card not in your hand: %w", sn.ErrValidation)
	case selectedCard.Rank.value() < roiRank.value() && cp.hasRank(roiRank):
		return nil, noCard, fmt.Errorf("must select a roi card that is not in your hand: %w", sn.ErrValidation)
	case selectedCard.Rank.value() < dameRank.value() && cp.hasRank(dameRank):
		return nil, noCard, fmt.Errorf("must select a dame card that is not in your hand: %w", sn.ErrValidation)
	case selectedCard.Rank.value() < cavalierRank.value() && cp.hasRank(cavalierRank):
		return nil, noCard, fmt.Errorf("must select a cavalier card that is not in your hand: %w", sn.ErrValidation)
	case selectedCard.Rank.value() < valetRank.value() && cp.hasRank(valetRank):
		return nil, noCard, fmt.Errorf("must select a valet card that is not in your hand: %w", sn.ErrValidation)
	case selectedCard.Rank.value() < tenRank.value() && cp.hasRank(tenRank):
		return nil, noCard, fmt.Errorf("must select a ten card that is not in your hand: %w", sn.ErrValidation)
	default:
		return cp, selectedCard, nil
	}
}
