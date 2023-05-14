package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
)

type space struct {
	rank rank
	kind kind
}

type kind string

const (
	noKind      kind = ""
	heartKind   kind = kind(hearts)
	clubKind    kind = kind(clubs)
	spadeKind   kind = kind(spades)
	diamondKind kind = kind(diamonds)
	trumpKind   kind = kind(trumps)
	trickKind   kind = "trick"
)

func (g *game) spacesFor(team []sn.PID) []space {
	ss := pie.Map(g.cardsFor(team), func(c card) space { return c.toSpace() })
	for i, win := range g.trickWonBy(team) {
		if win {
			ss = append(ss, space{toRank(i + 1), trickKind})
		}
	}
	return ss
}

func (g game) spacesNotOwnedBy(team []sn.PID) []space {
	_, notOwnedBy := pie.Diff(g.allSpaces(), g.spacesFor(team))
	sn.Debugf("notOwnedBy: %#v", notOwnedBy)
	return notOwnedBy
}

func (g game) allSpaces() []space {
	return pie.Keys(g.neighbors())
}

func (g game) neighbors() map[space][]space {
	if g.NumPlayers == 2 {
		return neighbors2()
	}
	return neighbors36()
}

func neighbors2() map[space][]space {
	return map[space][]space{

		// Row 1
		space{oneRank, trumpKind}: []space{
			space{oneRank, trickKind},
			space{cavalierRank, heartKind},
			space{fifteenRank, trickKind},
		},

		// Row 2
		space{fifteenRank, trickKind}: []space{
			space{fourteenRank, trickKind},
			space{valetRank, diamondKind},
			space{cavalierRank, heartKind},
			space{oneRank, trumpKind},
		},
		space{oneRank, trickKind}: []space{
			space{oneRank, trumpKind},
			space{cavalierRank, heartKind},
			space{valetRank, spadeKind},
			space{twoRank, trickKind},
		},

		// Row 3
		space{fourteenRank, trickKind}: []space{
			space{thirteenRank, trickKind},
			space{dameRank, heartKind},
			space{valetRank, diamondKind},
			space{fifteenRank, trickKind},
		},
		space{cavalierRank, heartKind}: []space{
			space{oneRank, trumpKind},
			space{fifteenRank, trickKind},
			space{valetRank, diamondKind},
			space{excuseRank, trumpKind},
			space{valetRank, spadeKind},
			space{oneRank, trickKind},
		},
		space{twoRank, trickKind}: []space{
			space{oneRank, trickKind},
			space{valetRank, spadeKind},
			space{dameRank, clubKind},
			space{threeRank, trickKind},
		},

		// Row 4
		space{thirteenRank, trickKind}: []space{
			space{twelveRank, trickKind},
			space{dameRank, heartKind},
			space{fourteenRank, trickKind},
		},
		space{valetRank, diamondKind}: []space{
			space{fifteenRank, trickKind},
			space{fourteenRank, trickKind},
			space{dameRank, heartKind},
			space{roiRank, spadeKind},
			space{excuseRank, trumpKind},
			space{cavalierRank, heartKind},
		},
		space{valetRank, spadeKind}: []space{
			space{oneRank, trickKind},
			space{cavalierRank, heartKind},
			space{excuseRank, trumpKind},
			space{roiRank, diamondKind},
			space{dameRank, clubKind},
			space{twoRank, trickKind},
		},
		space{threeRank, trickKind}: []space{
			space{twoRank, trickKind},
			space{dameRank, clubKind},
			space{fourRank, trickKind},
		},

		// Row 5
		space{dameRank, heartKind}: []space{
			space{fourteenRank, trickKind},
			space{thirteenRank, trickKind},
			space{twelveRank, trickKind},
			space{cavalierRank, clubKind},
			space{roiRank, spadeKind},
			space{valetRank, diamondKind},
		},
		space{excuseRank, trumpKind}: []space{
			space{cavalierRank, heartKind},
			space{valetRank, diamondKind},
			space{roiRank, spadeKind},
			space{sixteenRank, trickKind},
			space{roiRank, diamondKind},
			space{valetRank, spadeKind},
		},
		space{dameRank, clubKind}: []space{
			space{twoRank, trickKind},
			space{valetRank, spadeKind},
			space{roiRank, diamondKind},
			space{cavalierRank, spadeKind},
			space{fourRank, trickKind},
			space{threeRank, trickKind},
		},

		// Row 6
		space{twelveRank, trickKind}: []space{
			space{thirteenRank, trickKind},
			space{elevenRank, trickKind},
			space{cavalierRank, clubKind},
			space{dameRank, heartKind},
		},
		space{roiRank, spadeKind}: []space{
			space{valetRank, diamondKind},
			space{dameRank, heartKind},
			space{cavalierRank, clubKind},
			space{roiRank, heartKind},
			space{sixteenRank, trickKind},
			space{excuseRank, trumpKind},
		},
		space{roiRank, diamondKind}: []space{
			space{valetRank, spadeKind},
			space{excuseRank, trumpKind},
			space{sixteenRank, trickKind},
			space{roiRank, clubKind},
			space{cavalierRank, spadeKind},
			space{dameRank, clubKind},
		},
		space{fourRank, trickKind}: []space{
			space{threeRank, trickKind},
			space{dameRank, clubKind},
			space{cavalierRank, spadeKind},
			space{fiveRank, trickKind},
		},

		// Row 7
		space{cavalierRank, clubKind}: []space{
			space{dameRank, heartKind},
			space{twelveRank, trickKind},
			space{elevenRank, trickKind},
			space{dameRank, spadeKind},
			space{roiRank, heartKind},
			space{roiRank, spadeKind},
		},
		space{sixteenRank, trickKind}: []space{
			space{excuseRank, trumpKind},
			space{roiRank, spadeKind},
			space{roiRank, heartKind},
			space{twentyoneRank, trumpKind},
			space{roiRank, clubKind},
			space{roiRank, diamondKind},
		},
		space{cavalierRank, spadeKind}: []space{
			space{dameRank, clubKind},
			space{roiRank, diamondKind},
			space{roiRank, clubKind},
			space{dameRank, diamondKind},
			space{fiveRank, trickKind},
			space{fourRank, trickKind},
		},

		// Row 8
		space{elevenRank, trickKind}: []space{
			space{twelveRank, trickKind},
			space{thirteenRank, trumpKind},
			space{dameRank, spadeKind},
			space{cavalierRank, clubKind},
		},
		space{roiRank, heartKind}: []space{
			space{roiRank, spadeKind},
			space{cavalierRank, clubKind},
			space{dameRank, spadeKind},
			space{valetRank, clubKind},
			space{twentyoneRank, trumpKind},
			space{sixteenRank, trickKind},
		},
		space{roiRank, clubKind}: []space{
			space{roiRank, diamondKind},
			space{sixteenRank, trickKind},
			space{twentyoneRank, trumpKind},
			space{valetRank, heartKind},
			space{dameRank, diamondKind},
			space{cavalierRank, spadeKind},
		},
		space{fiveRank, trickKind}: []space{
			space{fourRank, trickKind},
			space{cavalierRank, spadeKind},
			space{dameRank, diamondKind},
			space{twoRank, trumpKind},
		},

		// Row 9
		space{dameRank, spadeKind}: []space{
			space{cavalierRank, clubKind},
			space{elevenRank, trickKind},
			space{threeRank, trumpKind},
			space{tenRank, trickKind},
			space{valetRank, clubKind},
			space{roiRank, heartKind},
		},
		space{twentyoneRank, trumpKind}: []space{
			space{sixteenRank, trickKind},
			space{roiRank, heartKind},
			space{valetRank, clubKind},
			space{cavalierRank, diamondKind},
			space{valetRank, heartKind},
			space{roiRank, diamondKind},
		},
		space{dameRank, diamondKind}: []space{
			space{cavalierRank, spadeKind},
			space{roiRank, clubKind},
			space{valetRank, heartKind},
			space{sixRank, trickKind},
			space{twoRank, trumpKind},
			space{fiveRank, trickKind},
		},

		// Row 10
		space{threeRank, trumpKind}: []space{
			space{elevenRank, trickKind},
			space{tenRank, trickKind},
			space{dameRank, spadeKind},
		},
		space{valetRank, clubKind}: []space{
			space{roiRank, heartKind},
			space{dameRank, spadeKind},
			space{tenRank, trickKind},
			space{nineRank, trickKind},
			space{cavalierRank, diamondKind},
			space{twentyoneRank, trumpKind},
		},
		space{valetRank, heartKind}: []space{
			space{roiRank, clubKind},
			space{twentyoneRank, trumpKind},
			space{cavalierRank, diamondKind},
			space{sevenRank, trickKind},
			space{sixRank, trickKind},
			space{dameRank, diamondKind},
		},
		space{twoRank, trumpKind}: []space{
			space{fiveRank, trickKind},
			space{dameRank, diamondKind},
			space{sixRank, trickKind},
		},

		// Row 11
		space{tenRank, trickKind}: []space{
			space{dameRank, spadeKind},
			space{threeRank, trumpKind},
			space{nineRank, trickKind},
			space{valetRank, clubKind},
		},
		space{cavalierRank, diamondKind}: []space{
			space{twentyoneRank, trumpKind},
			space{valetRank, clubKind},
			space{nineRank, trickKind},
			space{eightRank, trickKind},
			space{sevenRank, trickKind},
			space{valetRank, heartKind},
		},
		space{sixRank, trickKind}: []space{
			space{dameRank, diamondKind},
			space{valetRank, heartKind},
			space{sevenRank, trickKind},
			space{twoRank, trumpKind},
		},

		// Row 12
		space{nineRank, trickKind}: []space{
			space{valetRank, clubKind},
			space{tenRank, trickKind},
			space{eightRank, trickKind},
			space{cavalierRank, diamondKind},
		},
		space{sevenRank, trickKind}: []space{
			space{valetRank, heartKind},
			space{cavalierRank, diamondKind},
			space{eightRank, trickKind},
			space{sixRank, trickKind},
		},

		// Row 13
		space{eightRank, trickKind}: []space{
			space{cavalierRank, diamondKind},
			space{nineRank, trickKind},
			space{sevenRank, trickKind},
		},
	}
}
func neighbors36() map[space][]space {
	return map[space][]space{

		// Row 1
		space{oneRank, trumpKind}: []space{
			space{oneRank, trickKind},
			space{cavalierRank, heartKind},
			space{twelveRank, trickKind},
		},

		// Row 2
		space{twelveRank, trickKind}: []space{
			space{elevenRank, trickKind},
			space{valetRank, diamondKind},
			space{cavalierRank, heartKind},
			space{oneRank, trumpKind},
		},
		space{oneRank, trickKind}: []space{
			space{oneRank, trumpKind},
			space{cavalierRank, heartKind},
			space{valetRank, spadeKind},
			space{twoRank, trickKind},
		},

		// Row 3
		space{elevenRank, trickKind}: []space{
			space{sixRank, trumpKind},
			space{dameRank, heartKind},
			space{valetRank, diamondKind},
			space{twelveRank, trickKind},
		},
		space{cavalierRank, heartKind}: []space{
			space{oneRank, trumpKind},
			space{twelveRank, trickKind},
			space{valetRank, diamondKind},
			space{excuseRank, trumpKind},
			space{valetRank, spadeKind},
			space{oneRank, trickKind},
		},
		space{twoRank, trickKind}: []space{
			space{oneRank, trickKind},
			space{valetRank, spadeKind},
			space{dameRank, clubKind},
			space{twoRank, trumpKind},
		},

		// Row 4
		space{sixRank, trumpKind}: []space{
			space{tenRank, trickKind},
			space{dameRank, heartKind},
			space{elevenRank, trickKind},
		},
		space{valetRank, diamondKind}: []space{
			space{twelveRank, trickKind},
			space{elevenRank, trickKind},
			space{dameRank, heartKind},
			space{roiRank, spadeKind},
			space{excuseRank, trumpKind},
			space{cavalierRank, heartKind},
		},
		space{valetRank, spadeKind}: []space{
			space{oneRank, trickKind},
			space{cavalierRank, heartKind},
			space{excuseRank, trumpKind},
			space{roiRank, diamondKind},
			space{dameRank, clubKind},
			space{twoRank, trickKind},
		},
		space{twoRank, trumpKind}: []space{
			space{twoRank, trickKind},
			space{dameRank, clubKind},
			space{threeRank, trickKind},
		},

		// Row 5
		space{dameRank, heartKind}: []space{
			space{elevenRank, trickKind},
			space{sixRank, trumpKind},
			space{tenRank, trickKind},
			space{cavalierRank, clubKind},
			space{roiRank, spadeKind},
			space{valetRank, diamondKind},
		},
		space{excuseRank, trumpKind}: []space{
			space{cavalierRank, heartKind},
			space{valetRank, diamondKind},
			space{roiRank, spadeKind},
			space{thirteenRank, trickKind},
			space{roiRank, diamondKind},
			space{valetRank, spadeKind},
		},
		space{dameRank, clubKind}: []space{
			space{twoRank, trickKind},
			space{valetRank, spadeKind},
			space{roiRank, diamondKind},
			space{cavalierRank, spadeKind},
			space{threeRank, trickKind},
			space{twoRank, trumpKind},
		},

		// Row 6
		space{tenRank, trickKind}: []space{
			space{sixRank, trumpKind},
			space{nineRank, trickKind},
			space{cavalierRank, clubKind},
			space{dameRank, heartKind},
		},
		space{roiRank, spadeKind}: []space{
			space{valetRank, diamondKind},
			space{dameRank, heartKind},
			space{cavalierRank, clubKind},
			space{roiRank, heartKind},
			space{thirteenRank, trickKind},
			space{excuseRank, trumpKind},
		},
		space{roiRank, diamondKind}: []space{
			space{valetRank, spadeKind},
			space{excuseRank, trumpKind},
			space{thirteenRank, trickKind},
			space{roiRank, clubKind},
			space{cavalierRank, spadeKind},
			space{dameRank, clubKind},
		},
		space{threeRank, trickKind}: []space{
			space{twoRank, trumpKind},
			space{dameRank, clubKind},
			space{cavalierRank, spadeKind},
			space{fourRank, trickKind},
		},

		// Row 7
		space{cavalierRank, clubKind}: []space{
			space{dameRank, heartKind},
			space{tenRank, trickKind},
			space{nineRank, trickKind},
			space{dameRank, spadeKind},
			space{roiRank, heartKind},
			space{roiRank, spadeKind},
		},
		space{thirteenRank, trickKind}: []space{
			space{excuseRank, trumpKind},
			space{roiRank, spadeKind},
			space{roiRank, heartKind},
			space{twentyoneRank, trumpKind},
			space{roiRank, clubKind},
			space{roiRank, diamondKind},
		},
		space{cavalierRank, spadeKind}: []space{
			space{dameRank, clubKind},
			space{roiRank, diamondKind},
			space{roiRank, clubKind},
			space{dameRank, diamondKind},
			space{fourRank, trickKind},
			space{threeRank, trickKind},
		},

		// Row 9
		space{nineRank, trickKind}: []space{
			space{tenRank, trickKind},
			space{fiveRank, trumpKind},
			space{dameRank, spadeKind},
			space{cavalierRank, clubKind},
		},
		space{roiRank, heartKind}: []space{
			space{roiRank, spadeKind},
			space{cavalierRank, clubKind},
			space{dameRank, spadeKind},
			space{valetRank, clubKind},
			space{twentyoneRank, trumpKind},
			space{thirteenRank, trickKind},
		},
		space{roiRank, clubKind}: []space{
			space{roiRank, diamondKind},
			space{thirteenRank, trickKind},
			space{twentyoneRank, trumpKind},
			space{valetRank, heartKind},
			space{dameRank, diamondKind},
			space{cavalierRank, spadeKind},
		},
		space{fourRank, trickKind}: []space{
			space{threeRank, trickKind},
			space{cavalierRank, spadeKind},
			space{dameRank, diamondKind},
			space{threeRank, trumpKind},
		},

		// Row 10
		space{dameRank, spadeKind}: []space{
			space{cavalierRank, clubKind},
			space{nineRank, trickKind},
			space{fiveRank, trumpKind},
			space{eightRank, trickKind},
			space{valetRank, clubKind},
			space{roiRank, heartKind},
		},
		space{twentyoneRank, trumpKind}: []space{
			space{thirteenRank, trickKind},
			space{roiRank, heartKind},
			space{valetRank, clubKind},
			space{cavalierRank, diamondKind},
			space{valetRank, heartKind},
			space{roiRank, diamondKind},
		},
		space{dameRank, diamondKind}: []space{
			space{cavalierRank, spadeKind},
			space{roiRank, clubKind},
			space{valetRank, heartKind},
			space{fiveRank, trickKind},
			space{threeRank, trumpKind},
			space{fourRank, trickKind},
		},

		// Row 11
		space{fiveRank, trumpKind}: []space{
			space{nineRank, trickKind},
			space{eightRank, trickKind},
			space{dameRank, spadeKind},
		},
		space{valetRank, clubKind}: []space{
			space{roiRank, heartKind},
			space{dameRank, spadeKind},
			space{eightRank, trickKind},
			space{sevenRank, trickKind},
			space{cavalierRank, diamondKind},
			space{twentyoneRank, trumpKind},
		},
		space{valetRank, heartKind}: []space{
			space{roiRank, clubKind},
			space{twentyoneRank, trumpKind},
			space{cavalierRank, diamondKind},
			space{sixRank, trickKind},
			space{fiveRank, trickKind},
			space{dameRank, diamondKind},
		},
		space{threeRank, trumpKind}: []space{
			space{fourRank, trickKind},
			space{dameRank, diamondKind},
			space{fiveRank, trickKind},
		},

		// Row 12
		space{eightRank, trickKind}: []space{
			space{dameRank, spadeKind},
			space{fiveRank, trumpKind},
			space{sevenRank, trickKind},
			space{valetRank, clubKind},
		},
		space{cavalierRank, diamondKind}: []space{
			space{twentyoneRank, trumpKind},
			space{valetRank, clubKind},
			space{sevenRank, trickKind},
			space{fourRank, trumpKind},
			space{sixRank, trickKind},
			space{valetRank, heartKind},
		},
		space{fiveRank, trickKind}: []space{
			space{dameRank, diamondKind},
			space{valetRank, heartKind},
			space{sixRank, trickKind},
			space{threeRank, trumpKind},
		},

		// Row 13
		space{sevenRank, trickKind}: []space{
			space{valetRank, clubKind},
			space{eightRank, trickKind},
			space{fourRank, trumpKind},
			space{cavalierRank, diamondKind},
		},
		space{sixRank, trickKind}: []space{
			space{valetRank, heartKind},
			space{cavalierRank, diamondKind},
			space{fourRank, trumpKind},
			space{fiveRank, trickKind},
		},

		// Row 14
		space{fourRank, trumpKind}: []space{
			space{cavalierRank, diamondKind},
			space{sevenRank, trickKind},
			space{sixRank, trickKind},
		},
	}
}
