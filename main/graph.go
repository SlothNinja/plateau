package main

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func (g game) graphFor(t spaceTest) boardGraph {
	graph := newBoardGraph()

	// add all spaces to graph
	pie.Each(allSpaces(), func(s space) {
		u := graph.UndirectedGraph.NewNode()
		uid := u.ID()
		u = node{space: s, id: uid}
		graph.UndirectedGraph.AddNode(u)
		graph.ids[s] = uid
	})

	graph.addEdges(g, t)
	return graph
}

func allSpaces() []space {
	return pie.Keys(neighbors)
}

type boardGraph struct {
	ids map[space]int64
	*simple.UndirectedGraph
}

type spaceTest func(space) bool

func (graph boardGraph) addEdges(g game, t spaceTest) {
	pie.Each(g.spacesFor(g.declarersTeam), func(s space) {
		ns, ok := neighbors[s]
		// if no neighbors then space not part of board
		if !ok {
			return
		}

		// should not occur since all spaces already added to graph
		if _, exists := graph.ids[s]; !exists {
			sn.Warningf("space %q not in graph", s)
			return
		}

		pie.Each(ns, func(n space) {
			sn.Debugf("space: %#v", n)
			if t(n) {
				sn.Debugf("added: %#v", n)
				sid, ok := graph.ids[s]
				if !ok {
					sn.Warningf("space %q not in graph", s)
					return
				}
				u := graph.Node(sid)
				if u == nil {
					sn.Warningf("node %q not in graph", u)
					return
				}
				nid, ok := graph.ids[n]
				if !ok {
					sn.Warningf("space %q not in graph", n)
					return
				}

				v := graph.Node(nid)
				if v == nil {
					sn.Warningf("node %q not in graph", v)
					return
				}
				graph.SetEdge(simple.Edge{F: u, T: v})
			}
		})
	})
}

func own(g game, team []sn.PID) spaceTest {
	return func(s space) bool {
		return pie.Contains(g.spacesFor(team), s)
	}
}

func mayOwn(g game, team []sn.PID) spaceTest {
	return func(s space) bool {
		return !own(g, g.otherTeam(team))(s)
	}
}

func newBoardGraph() boardGraph {
	return boardGraph{
		ids:             make(map[space]int64),
		UndirectedGraph: simple.NewUndirectedGraph(),
	}
}

type node struct {
	space space
	id    int64
}

// implement graph.Node interface
func (n node) ID() int64 { return n.id }

func (g *game) cardsFor(team []sn.PID) []card {
	var cards []card
	pie.Each(g.tricksFor(team), func(t trick) {
		cards = append(cards, pie.Map(t.cards, func(c card) card {
			c.playedBy = sn.NoPID
			return c
		})...)
	})
	return cards
}

func (g *game) spacesFor(team []sn.PID) []space {
	ss := pie.Map(g.cardsFor(team), func(c card) space { return c.toSpace() })
	for i, win := range g.trickWonBy(team) {
		if win {
			ss = append(ss, space{toRank(i + 1), trickKind})
		}
	}
	return ss
}

func (g *game) tricksFor(team []sn.PID) []trick {
	return pie.Filter(g.tricks, func(t trick) bool { return pie.Contains(team, t.wonBy) })
}

func (g *game) trickWonBy(team []sn.PID) (won []bool) {
	pie.Each(g.tricks, func(t trick) {
		won = append(won, pie.Contains(team, t.wonBy))
	})
	return won
}

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

func (c card) toSpace() space {
	return space{c.rank, kind(c.suit)}
}

var neighbors map[space][]space = map[space][]space{

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

func bridge(paths path.AllShortest, s1, s2, s3, s4, s5, s6 []node) bool {
	return connected2(paths, s1, s4) || connected2(paths, s2, s5) || connected2(paths, s3, s6)
}

func y(paths path.AllShortest, s1, s2, s3, s4, s5, s6 []node) bool {
	return connected3(paths, s1, s3, s5) || connected3(paths, s2, s4, s6)
}

func fork(paths path.AllShortest, s1, s2, s3, s4, s5, s6 []node) bool {
	return bridge(paths, s1, s2, s3, s4, s5, s6) && y(paths, s1, s2, s3, s4, s5, s6)
}

func fiveSides(paths path.AllShortest, s1, s2, s3, s4, s5, s6 []node) bool {
	return connected5(paths, s1, s2, s3, s4, s5) ||
		connected5(paths, s1, s2, s3, s4, s6) ||
		connected5(paths, s1, s2, s3, s5, s6) ||
		connected5(paths, s1, s2, s4, s5, s6) ||
		connected5(paths, s1, s3, s4, s5, s6) ||
		connected5(paths, s2, s3, s4, s5, s6)
}

func sixSides(paths path.AllShortest, s1, s2, s3, s4, s5, s6 []node) bool {
	return pie.Any(s1, func(n1 node) bool {
		return pie.Any(s2, func(n2 node) bool {
			p, _, _ := paths.Between(n1.ID(), n2.ID())
			return p != nil
		}) &&
			pie.Any(s3, func(n3 node) bool {
				p, _, _ := paths.Between(n1.ID(), n3.ID())
				return p != nil
			}) &&
			pie.Any(s4, func(n4 node) bool {
				p, _, _ := paths.Between(n1.ID(), n4.ID())
				return p != nil
			}) &&
			pie.Any(s5, func(n5 node) bool {
				p, _, _ := paths.Between(n1.ID(), n5.ID())
				return p != nil
			}) &&
			pie.Any(s6, func(n6 node) bool {
				p, _, _ := paths.Between(n1.ID(), n6.ID())
				return p != nil
			})
	})
}

func connected2(paths path.AllShortest, s1, s2 []node) bool {
	return pie.Any(s1, func(n1 node) bool {
		return pie.Any(s2, func(n2 node) bool {
			p, _, _ := paths.Between(n1.ID(), n2.ID())
			return p != nil
		})
	})
}

func connected3(paths path.AllShortest, s1, s2, s3 []node) bool {
	return pie.Any(s1, func(n1 node) bool {
		return pie.Any(s2, func(n2 node) bool {
			p, _, _ := paths.Between(n1.ID(), n2.ID())
			return p != nil
		}) &&
			pie.Any(s3, func(n3 node) bool {
				p, _, _ := paths.Between(n1.ID(), n3.ID())
				return p != nil
			})
	})
}

func connected5(paths path.AllShortest, s1, s2, s3, s4, s5 []node) bool {
	return pie.Any(s1, func(n1 node) bool {
		return pie.Any(s2, func(n2 node) bool {
			p, _, _ := paths.Between(n1.ID(), n2.ID())
			return p != nil
		}) &&
			pie.Any(s3, func(n3 node) bool {
				p, _, _ := paths.Between(n1.ID(), n3.ID())
				return p != nil
			}) &&
			pie.Any(s4, func(n4 node) bool {
				p, _, _ := paths.Between(n1.ID(), n4.ID())
				return p != nil
			}) &&
			pie.Any(s5, func(n5 node) bool {
				p, _, _ := paths.Between(n1.ID(), n5.ID())
				return p != nil
			})
	})
}

func side1(graph boardGraph) (nodes []node) {
	pie.Each([]space{
		space{oneRank, trumpKind},
		space{oneRank, trickKind},
		space{twoRank, trickKind},
		space{twoRank, trumpKind},
	}, func(s space) { nodes = append(nodes, node{s, graph.Node(graph.ids[s]).ID()}) })
	return nodes
}

func side2(graph boardGraph) (nodes []node) {
	pie.Each([]space{
		space{twoRank, trumpKind},
		space{threeRank, trickKind},
		space{fourRank, trickKind},
		space{threeRank, trumpKind},
	}, func(s space) { nodes = append(nodes, node{s, graph.Node(graph.ids[s]).ID()}) })
	return nodes
}

func side3(graph boardGraph) (nodes []node) {
	pie.Each([]space{
		space{threeRank, trumpKind},
		space{fiveRank, trickKind},
		space{sixRank, trickKind},
		space{fourRank, trumpKind},
	}, func(s space) { nodes = append(nodes, node{s, graph.Node(graph.ids[s]).ID()}) })
	return nodes
}

func side4(graph boardGraph) (nodes []node) {
	pie.Each([]space{
		space{fourRank, trumpKind},
		space{sevenRank, trickKind},
		space{eightRank, trickKind},
		space{fiveRank, trumpKind},
	}, func(s space) { nodes = append(nodes, node{s, graph.Node(graph.ids[s]).ID()}) })
	return nodes
}

func side5(graph boardGraph) (nodes []node) {
	pie.Each([]space{
		space{fiveRank, trumpKind},
		space{nineRank, trickKind},
		space{tenRank, trickKind},
		space{sixRank, trumpKind},
	}, func(s space) { nodes = append(nodes, node{s, graph.Node(graph.ids[s]).ID()}) })
	return nodes
}

func side6(graph boardGraph) (nodes []node) {
	pie.Each([]space{
		space{sixRank, trumpKind},
		space{elevenRank, trickKind},
		space{twelveRank, trickKind},
		space{oneRank, trumpKind},
	}, func(s space) { nodes = append(nodes, node{s, graph.Node(graph.ids[s]).ID()}) })
	return nodes
}
