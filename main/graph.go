package main

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

		ns, ok := neighbors[s]
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

func bridge(graph boardGraph, paths path.AllShortest) (path []node, found bool) {
	path, found = connected(paths, graph.side(1), graph.side(4))
	if found {
		return path, found
	}
	path, found = connected(paths, graph.side(2), graph.side(5))
	if found {
		return path, found
	}
	return connected(paths, graph.side(3), graph.side(6))
}

func y(graph boardGraph, paths path.AllShortest) (path []node, found bool) {
	sn.Debugf(msgEnter)
	defer sn.Debugf(msgExit)

	path, found = connected(paths, graph.side(1), graph.side(3), graph.side(5))
	if found {
		return path, true
	}
	return connected(paths, graph.side(2), graph.side(4), graph.side(6))
}

func fork(graph boardGraph, paths path.AllShortest) (path []node, found bool) {
	path1, found1 := bridge(graph, paths)
	if !found1 {
		return nil, false
	}
	path2, found2 := y(graph, paths)
	if !found2 {
		return nil, false
	}
	return append(path1, path2...), true
}

func fiveSides(graph boardGraph, paths path.AllShortest) ([]node, bool) {
	s1, s2, s3, s4, s5, s6 := graph.side(1), graph.side(2), graph.side(3), graph.side(4), graph.side(5), graph.side(6)
	path, found := connected(paths, s1, s2, s3, s4, s5)
	if found {
		return path, true
	}
	path, found = connected(paths, s1, s2, s3, s4, s6)
	if found {
		return path, true
	}
	path, found = connected(paths, s1, s2, s3, s5, s6)
	if found {
		return path, true
	}
	path, found = connected(paths, s1, s2, s4, s5, s6)
	if found {
		return path, true
	}
	path, found = connected(paths, s1, s3, s4, s5, s6)
	if found {
		return path, true
	}
	path, found = connected(paths, s2, s3, s4, s5, s6)
	if found {
		return path, true
	}
	return nil, false
}

func sixSides(graph boardGraph, paths path.AllShortest) ([]node, bool) {
	s1, s2, s3, s4, s5, s6 := graph.side(1), graph.side(2), graph.side(3), graph.side(4), graph.side(5), graph.side(6)
	return connected(paths, s1, s2, s3, s4, s5, s6)
}

func connected(paths path.AllShortest, ss ...[]node) (path []node, found bool) {
	found = pie.Any(pie.First(ss), func(n0 node) bool {
		return pie.All(pie.DropTop(ss, 1), func(ss1 []node) bool {
			return pie.Any(ss1, func(n1 node) bool {
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

func (graph boardGraph) side(s int) (nodes []node) {
	pie.Each(map[int][]space{
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
	}[s], func(s space) {
		if nid, exists := graph.ids[s]; exists {
			nodes = append(nodes, graph.Node(nid).(node))
		}
	})
	return nodes
}

func toNodes(path []graph.Node) []node {
	nodes := make([]node, len(path))
	for i := range path {
		nodes[i] = path[i].(node)
	}
	return nodes
}
