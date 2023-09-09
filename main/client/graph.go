package client

import (
	"github.com/SlothNinja/sn/v3"
	"github.com/elliotchance/pie/v2"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

func (g game) graphFor(ss []space) boardGraph {
	graph := newBoardGraph()

	// add all spaces to graph
	for _, s := range ss {
		u := graph.UndirectedGraph.NewNode()
		uid := u.ID()
		u = node{Space: s, UID: uid}
		graph.UndirectedGraph.AddNode(u)
		graph.ids[s] = uid

		ns, ok := g.neighbors()[s]
		if !ok {
			continue
		}
		for _, n := range ns {
			if _, exists := graph.ids[s]; !exists {
				continue
			}
			if pie.Contains(ss, n) {
				nid, ok := graph.ids[n]
				if !ok {
					continue
				}

				v := graph.Node(nid)
				if v == nil {
					continue
				}
				graph.SetEdge(simple.Edge{F: u, T: v})
			}
		}
	}
	return graph

}

type boardGraph struct {
	ids map[space]int64
	*simple.UndirectedGraph
}

func newBoardGraph() boardGraph {
	return boardGraph{
		ids:             make(map[space]int64),
		UndirectedGraph: simple.NewUndirectedGraph(),
	}
}

type node struct {
	Space space
	UID   int64
}

// implement graph.Node interface
func (n node) ID() int64 { return n.UID }

func (g game) bridge(graph boardGraph, paths path.AllShortest) ([]node, handResult) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	path, connected := isConnected(paths, g.side(graph, 1), g.side(graph, 4))
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, g.side(graph, 2), g.side(graph, 5))
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, g.side(graph, 3), g.side(graph, 6))
	if connected {
		return path, dSuccess
	}
	return nil, dFail
}

func (g game) y(graph boardGraph, paths path.AllShortest) ([]node, handResult) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	path, connected := isConnected(paths, g.side(graph, 1), g.side(graph, 3), g.side(graph, 5))
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, g.side(graph, 2), g.side(graph, 4), g.side(graph, 6))
	if connected {
		return path, dSuccess
	}
	return nil, dFail
}

func (g game) fork(graph boardGraph, paths path.AllShortest) (path []node, result handResult) {
	path1, result1 := g.bridge(graph, paths)
	if result1 == dFail {
		return nil, result1
	}
	path2, result2 := g.y(graph, paths)
	if result2 == dFail {
		return nil, result2
	}
	return append(path1, path2...), dSuccess
}

func (g game) fiveSides(graph boardGraph, paths path.AllShortest) ([]node, handResult) {
	s1, s2, s3, s4, s5, s6 := g.side(graph, 1), g.side(graph, 2), g.side(graph, 3), g.side(graph, 4), g.side(graph, 5), g.side(graph, 6)

	path, connected := isConnected(paths, s1, s2, s3, s4, s5)
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, s1, s2, s3, s4, s6)
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, s1, s2, s3, s5, s6)
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, s1, s2, s4, s5, s6)
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, s1, s3, s4, s5, s6)
	if connected {
		return path, dSuccess
	}

	path, connected = isConnected(paths, s2, s3, s4, s5, s6)
	if connected {
		return path, dSuccess
	}
	return nil, dFail
}

func (g game) sixSides(graph boardGraph, paths path.AllShortest) ([]node, handResult) {
	s1, s2, s3, s4, s5, s6 := g.side(graph, 1), g.side(graph, 2), g.side(graph, 3), g.side(graph, 4), g.side(graph, 5), g.side(graph, 6)

	path, connected := isConnected(paths, s1, s2, s3, s4, s5, s6)
	if connected {
		return path, dSuccess
	}
	return nil, dFail
}

func isConnected(paths path.AllShortest, sides ...[]node) (path []node, result bool) {
	found := pie.Any(pie.First(sides), func(n0 node) bool {
		return pie.All(pie.DropTop(sides, 1), func(side []node) bool {
			return pie.Any(side, func(n1 node) bool {
				p, _, _ := paths.Between(n0.ID(), n1.ID())
				if p != nil {
					path = append(path, toNodes(p)...)
					return true
				}
				return false
			})
		})
	})
	if found {
		return path, true
	}
	return nil, false
}

func (g game) side(graph boardGraph, s int) (nodes []node) {
	var sideSpaces map[int][]space
	if g.NumPlayers == 2 {
		sideSpaces = sides2()
	} else {
		sideSpaces = sides36()
	}
	pie.Each(sideSpaces[s], func(s space) {
		if nid, exists := graph.ids[s]; exists {
			nodes = append(nodes, graph.Node(nid).(node))
		}
	})
	return nodes
}

func sides2() map[int][]space {
	return map[int][]space{
		1: []space{
			space{oneRank, trumpKind},
			space{oneRank, trickKind},
			space{twoRank, trickKind},
			space{threeRank, trickKind},
		},
		2: []space{
			space{threeRank, trickKind},
			space{fourRank, trickKind},
			space{fiveRank, trickKind},
			space{twoRank, trumpKind},
		},
		3: []space{
			space{twoRank, trumpKind},
			space{sixRank, trickKind},
			space{sevenRank, trickKind},
			space{eightRank, trickKind},
		},
		4: []space{
			space{eightRank, trickKind},
			space{nineRank, trickKind},
			space{tenRank, trickKind},
			space{threeRank, trumpKind},
		},
		5: []space{
			space{threeRank, trumpKind},
			space{elevenRank, trickKind},
			space{twelveRank, trickKind},
			space{thirteenRank, trickKind},
		},
		6: []space{
			space{thirteenRank, trickKind},
			space{fourteenRank, trickKind},
			space{fifteenRank, trickKind},
			space{oneRank, trumpKind},
		},
	}
}

func sides36() map[int][]space {
	return map[int][]space{
		1: []space{
			space{oneRank, trumpKind},
			space{oneRank, trickKind},
			space{twoRank, trickKind},
			space{twoRank, trumpKind},
		},
		2: []space{
			space{twoRank, trumpKind},
			space{threeRank, trickKind},
			space{fourRank, trickKind},
			space{threeRank, trumpKind},
		},
		3: []space{
			space{threeRank, trumpKind},
			space{fiveRank, trickKind},
			space{sixRank, trickKind},
			space{fourRank, trumpKind},
		},
		4: []space{
			space{fourRank, trumpKind},
			space{sevenRank, trickKind},
			space{eightRank, trickKind},
			space{fiveRank, trumpKind},
		},
		5: []space{
			space{fiveRank, trumpKind},
			space{nineRank, trickKind},
			space{tenRank, trickKind},
			space{sixRank, trumpKind},
		},
		6: []space{
			space{sixRank, trumpKind},
			space{elevenRank, trickKind},
			space{twelveRank, trickKind},
			space{oneRank, trumpKind},
		},
	}
}

func toNodes(path []graph.Node) []node {
	nodes := make([]node, len(path))
	for i := range path {
		nodes[i] = path[i].(node)
	}
	return nodes
}
