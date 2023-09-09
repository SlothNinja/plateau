package client

import (
	"testing"

	"github.com/SlothNinja/sn/v3"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/path"
)

func TestConnected2(t *testing.T) {

	var g game
	g.DeclarersTeam = []sn.PID{1}
	g.Tricks = []trick{
		trick{[]card{
			card{oneRank, trumps, 1, true},
			card{cavalierRank, hearts, 2, true},
			card{excuseRank, trumps, 3, true},
			card{roiRank, spades, 4, true},
			card{roiRank, hearts, 5, true},
			card{twentyoneRank, trumps, 6, true},
		}, 1},
		trick{[]card{
			card{cavalierRank, diamonds, 1, true},
			card{fourRank, trumps, 2, true},
			card{excuseRank, trumps, 3, true},
			card{tenRank, spades, 4, true},
			card{tenRank, hearts, 5, true},
			card{tenRank, clubs, 6, true},
		}, 1},
	}

	spaces := g.spacesFor(g.DeclarersTeam)
	assert.NotNil(t, spaces, "spaces should not be nil")

	graph := g.graphFor(spaces)
	paths := path.DijkstraAllPaths(graph)
	s1 := g.side(graph, 1)
	s2 := g.side(graph, 2)
	s3 := g.side(graph, 3)
	s4 := g.side(graph, 4)
	s5 := g.side(graph, 5)
	s6 := g.side(graph, 6)

	path, connected := isConnected(paths, s1, s2)
	assert.Falsef(t, connected, "path should not be found between side1 and side2")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3)
	assert.Truef(t, connected, "path should be found between side1 and side3")
	assert.NotNilf(t, path, "path not should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s4)
	assert.Truef(t, connected, "path should be found between side1 and side4")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s5)
	assert.Falsef(t, connected, "path should not be found between side1 and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s6)
	assert.Truef(t, connected, "path should not be found between side1 and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)
}

func TestConnected3(t *testing.T) {

	var g game
	g.DeclarersTeam = []sn.PID{1}
	g.Tricks = []trick{
		trick{[]card{
			card{oneRank, trumps, 1, true},
			card{cavalierRank, hearts, 2, true},
			card{excuseRank, trumps, 3, true},
			card{roiRank, spades, 4, true},
			card{roiRank, hearts, 5, true},
			card{twentyoneRank, trumps, 6, true},
		}, 1},
		trick{[]card{
			card{cavalierRank, diamonds, 1, true},
			card{fourRank, trumps, 2, true},
			card{excuseRank, trumps, 3, true},
			card{tenRank, spades, 4, true},
			card{tenRank, hearts, 5, true},
			card{tenRank, clubs, 6, true},
		}, 1},
	}

	spaces := g.spacesFor(g.DeclarersTeam)
	assert.NotNil(t, spaces, "spaces should not be nil")

	graph := g.graphFor(spaces)
	paths := path.DijkstraAllPaths(graph)
	s1 := g.side(graph, 1)
	s2 := g.side(graph, 2)
	s3 := g.side(graph, 3)
	s4 := g.side(graph, 4)
	s5 := g.side(graph, 5)
	s6 := g.side(graph, 6)

	path, connected := isConnected(paths, s1, s2, s3)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side3")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s2, s4)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side4")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s2, s5)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s2, s6)
	assert.Falsef(t, connected, "path should not be found between side1, side2, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3, s4)
	assert.Truef(t, connected, "path should be found between side1, side3, and side4")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3, s5)
	assert.Falsef(t, connected, "path should not be found between side1, side3, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s3, s6)
	assert.Truef(t, connected, "path should be found between side1, side3, and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s4, s5)
	assert.Falsef(t, connected, "path should not be found between side1, side4, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s1, s4, s6)
	assert.Truef(t, connected, "path should be found between side1, side4, and side6")
	assert.NotNilf(t, path, "path should not be nil: %#v", path)

	path, connected = isConnected(paths, s1, s5, s6)
	assert.Falsef(t, connected, "path should not be found between side1, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s2, s3, s4)
	assert.Falsef(t, connected, "path should not be found between side2, side3, and side4")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s2, s3, s5)
	assert.Falsef(t, connected, "path should not be found between side2, side3, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s2, s3, s6)
	assert.Falsef(t, connected, "path should not be found between side2, side3, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s3, s4, s5)
	assert.Falsef(t, connected, "path should not be found between side3, side4, and side5")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s3, s5, s6)
	assert.Falsef(t, connected, "path should not be found between side3, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)

	path, connected = isConnected(paths, s4, s5, s6)
	assert.Falsef(t, connected, "path should not be found between side4, side5, and side6")
	assert.Nilf(t, path, "path should be nil: %#v", path)
}
