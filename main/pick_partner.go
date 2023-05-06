package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

func (g *game) startPickPartner() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = pickPartnerPhase
	switch b := g.lastBid(); b.Teams {
	case "duoBid":
		return g.declarer()
	case "trioBid":
		return g.declarer()
	default:
		return nil
	}
}

func (g game) otherTeam(pids1 []sn.PID) []sn.PID {
	return pie.FilterNot(pidsFor(g.Players), func(pid2 sn.PID) bool {
		return pie.Any(pids1, func(pid1 sn.PID) bool { return pid1 == pid2 })
	})
}
