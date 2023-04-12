package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

func (g *game) startPickPartner() *player {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	g.Phase = pickPartnerPhase
	switch g.lastBid().teams {
	case "duoBid":
	case "trioBid":
	default:
		g.appendEntry(message{"template": "no-partner-picked"})
	}
	return g.startIncObjective()
}

func (g *game) otherTeam(pids1 []sn.PID) []sn.PID {
	return pie.FilterNot(pidsFor(g.players), func(pid2 sn.PID) bool {
		return pie.Any(pids1, func(pid1 sn.PID) bool { return pid1 == pid2 })
	})
}
